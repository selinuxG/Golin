package run

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// privShell 用"一次提权"换一个常驻 root shell：
// 内置模板采集的七十多条命令全部在这个 shell 里执行，
// 不再每条命令都新建 sudo/su 会话（原来每台 20~40s，现在 1~2s，
// 也不会在 /var/log/secure 里堆七十多条 su 记录）。
type privShell struct {
	mu     sync.Mutex
	sess   *ssh.Session
	in     io.WriteCloser
	out    *bufio.Reader
	seq    int
	closed bool
}

// newPrivShell 建立提权 shell：
//
//	su  模式：su -l <user>        （口令先喂给 stdin）
//	sudo模式：sudo -S -p '' bash  （口令先喂给 stdin）
//
// 建完后用 id -u 探针确认真的到 root，否则返回错误由调用方回退。
func newPrivShell(client *ssh.Client, mode string, user string, password string) (*privShell, error) {
	sess, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	in, err := sess.StdinPipe()
	if err != nil {
		sess.Close()
		return nil, err
	}
	stdout, err := sess.StdoutPipe()
	if err != nil {
		sess.Close()
		return nil, err
	}
	var start string
	switch mode {
	case "su":
		start = "su -l " + user
	case "sudo":
		start = "sudo -S -p '' bash"
	case "sudo-nopass":
		start = "sudo -n bash"
	case "none":
		start = "bash"
	default:
		sess.Close()
		return nil, fmt.Errorf("当前没有提权方式")
	}
	if err := sess.Start(start); err != nil {
		sess.Close()
		return nil, err
	}
	p := &privShell{sess: sess, in: in, out: bufio.NewReaderSize(stdout, 64*1024)}
	if password != "" {
		if _, err := io.WriteString(p.in, password+"\n"); err != nil {
			p.Close()
			return nil, err
		}
	}
	//探针：确认已经拿到 root（su/sudo 可能先打 motd/"Last failed login" 横幅，所以用标记取数）
	out, err := p.run(`echo GLNUID$(id -u)GLNEND`)
	if err != nil {
		p.Close()
		return nil, err
	}
	got := ""
	if i := strings.Index(out, "GLNUID"); i >= 0 {
		rest := out[i+len("GLNUID"):]
		if j := strings.Index(rest, "GLNEND"); j >= 0 {
			got = strings.TrimSpace(rest[:j])
		}
	}
	if mode == "none" {
		//不提权（登录用户就是 root）：只确认 shell 可用
		if got == "" {
			p.Close()
			return nil, fmt.Errorf("shell 不可用")
		}
	} else if got != "0" {
		p.Close()
		return nil, fmt.Errorf("提权失败(id -u=%q)", got)
	}
	return p, nil
}

// run 在这个 root shell 里执行一条命令，用哨兵行切分输出
func (p *privShell) run(cmd string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return "", fmt.Errorf("提权shell已关闭")
	}
	p.seq++
	tag := fmt.Sprintf("__GLN%d__", p.seq)
	//stderr 合并进 stdout，命令本身用 {} 包住避免 && / ; 影响哨兵
	line := "{ " + cmd + " ; } </dev/null 2>&1; echo " + tag + "E$?\n"
	if _, err := io.WriteString(p.in, line); err != nil {
		return "", err
	}
	type result struct {
		out string
		err error
	}
	ch := make(chan result, 1)
	go func() {
		var b strings.Builder
		for {
			l, err := p.out.ReadString('\n')
			if err != nil {
				ch <- result{b.String(), err}
				return
			}
			if strings.HasPrefix(l, tag+"E") {
				ch <- result{b.String(), nil}
				return
			}
			b.WriteString(l)
		}
	}()
	select {
	case r := <-ch:
		return r.out, r.err
	case <-time.After(120 * time.Second):
		//持锁状态下直接关，避免再走 Close() 造成死锁
		p.closed = true
		_ = p.in.Close()
		_ = p.sess.Close()
		return "", fmt.Errorf("命令超时")
	}
}

// Close 关掉提权 shell
func (p *privShell) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return
	}
	p.closed = true
	_, _ = io.WriteString(p.in, "exit\n")
	_ = p.in.Close()
	_ = p.sess.Close()
}
