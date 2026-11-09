package run

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

// SecCheck 一条核查结论
// Item 就是"检查的配置/命令"，Result 是依据（尽量直接引用已采集的数据），Judge 是判定。
// 除了 sshd -T 与时间同步这两项，其余全部由已采集到的字段计算，不再重复执行命令。
type SecCheck struct {
	Item   string //标题：检查项（短标题）
	Note   string //小子说明：检查的配置文件/命令（显示在标题下方）
	Result string //依据
	Judge  string //符合 / 部分符合 / 不符合 / 不适用 / 需人工确认 / 采集失败
}

// summaryNotes 标题 -> 小子说明（正文里的小字，不放在标题里）
var summaryNotes = map[string]string{
	"空口令账号":       "/etc/shadow：按 passwd -S 的密码状态（NP）统计",
	"账号唯一性":       "/etc/passwd：检查 UID 是否重复",
	"默认账号":        "/etc/passwd：adm、lp、games 等默认账号是否可登录",
	"口令复杂度策略":     "/etc/security/pwquality.conf 生效项 + PAM 是否引用 pam_pwquality",
	"登录失败处理":      "/etc/pam.d/system-auth、password-auth、sshd：pam_tally2 / pam_faillock、deny 阈值",
	"会话超时 TMOUT":  "/etc/profile、/etc/bashrc、/etc/profile.d",
	"日志外发":        "/etc/rsyslog.conf：是否配置 @@ 转发到远端日志服务器",
	"日志保留与切割":     "/etc/logrotate.conf、/etc/logrotate.d/syslog",
	"审计规则":        "auditctl -l：内核审计规则条数",
	"高危服务端口":      "ss -lnt：检查 21/23/69/512/513/514 端口",
	"防病毒与主机防护":    "ps 进程列表：匹配 clamd、safedog、aegis、edr、qax 等",
	"sudo 授权":     "/etc/sudoers、/etc/sudoers.d",
	"内核与系统版本":     "uname -r、/etc/*release",
	"SSH 生效配置与算法": "sshd -T：permitrootlogin、ciphers、macs、kexalgorithms 等",
	"时间同步":        "timedatectl / chronyc tracking / ntpq -p",
}

// 系统默认账号（等保要求改名或删除）
var defaultAccounts = map[string]bool{
	"adm": true, "lp": true, "games": true, "ftp": true, "news": true,
	"uucp": true, "operator": true, "sync": true, "shutdown": true, "halt": true,
}

// accountKind 账号类型：root / 默认账号（adm、lp、games 等）/ 系统账号（UID<1000）/ 普通账号
func accountKind(name string, uid string) string {
	if name == "root" {
		return "root"
	}
	if defaultAccounts[name] {
		return "默认账号"
	}
	u := strings.TrimSpace(uid)
	if u != "" {
		n, ok := 0, true
		for _, c := range u {
			if c < '0' || c > '9' {
				ok = false
				break
			}
			n = n*10 + int(c-'0')
		}
		if ok && n < 1000 {
			return "系统账号"
		}
	}
	return "普通账号"
}

// buildSummary 用已采集的数据生成核查结论；c 为 nil（本机模式）时跳过需要新命令的两项
func buildSummary(d *Data, c *ssh.Client) (out []SecCheck) {
	add := func(item string, fn func() SecCheck) {
		defer func() {
			if r := recover(); r != nil {
				out = append(out, SecCheck{Item: item, Result: fmt.Sprintf("计算异常：%v", r), Judge: "采集失败"})
			}
		}()
		sc := fn()
		if sc.Item == "" {
			sc.Item = item
		}
		sc.Note = summaryNotes[sc.Item]
		out = append(out, sc)
	}

	//1. 空口令账号 —— 直接看用户表里的密码状态（passwd -S 的 NP）
	add("空口令账号", func() SecCheck {
		if len(d.User) == 0 {
			return SecCheck{Result: "用户信息未采集到（需 root 权限读取 /etc/passwd、/etc/shadow）", Judge: "采集失败"}
		}
		var np []string
		for _, u := range d.User {
			if strings.EqualFold(strings.TrimSpace(u.Passwd), "NP") {
				np = append(np, u.Name)
			}
		}
		if len(np) == 0 {
			return SecCheck{Result: fmt.Sprintf("用户信息表共 %d 个账号，密码状态无 NP（空口令）", len(d.User)), Judge: "符合"}
		}
		return SecCheck{Result: "空口令账号：" + strings.Join(np, "、"), Judge: "不符合"}
	})

	//2. 重复 UID —— 用用户表 UID 统计
	add("账号唯一性", func() SecCheck {
		if len(d.User) == 0 {
			return SecCheck{Result: "用户信息未采集到", Judge: "采集失败"}
		}
		cnt := make(map[string][]string)
		for _, u := range d.User {
			uid := strings.TrimSpace(u.Uid)
			if uid == "" {
				continue
			}
			cnt[uid] = append(cnt[uid], u.Name)
		}
		var dup []string
		for uid, names := range cnt {
			if len(names) > 1 {
				dup = append(dup, fmt.Sprintf("UID %s：%s", uid, strings.Join(names, "、")))
			}
		}
		if len(dup) == 0 {
			return SecCheck{Result: fmt.Sprintf("用户信息表共 %d 个账号，无重复 UID", len(d.User)), Judge: "符合"}
		}
		return SecCheck{Result: "重复 UID：\n" + strings.Join(dup, "\n"), Judge: "不符合"}
	})

	//3. 默认账号是否可以登录 —— 用用户表里的默认账号标记 + 可登录列
	add("默认账号", func() SecCheck {
		if len(d.User) == 0 {
			return SecCheck{Result: "用户信息未采集到", Judge: "采集失败"}
		}
		var have, login []string
		for _, u := range d.User {
			if u.Kind != "默认账号" {
				continue
			}
			have = append(have, u.Name)
		}
		//可登录的默认/系统账号都算风险（UID<1000 且 shell 是交互式）
		for _, u := range d.User {
			if (u.Kind == "默认账号" || u.Kind == "系统账号") && u.Login {
				login = append(login, u.Name+"("+u.Bash+")")
			}
		}
		if len(have) == 0 && len(login) == 0 {
			return SecCheck{Result: "未发现系统默认账号（adm/lp/games 等均已删除），也无其他可登录的系统账号", Judge: "符合"}
		}
		if len(login) > 0 {
			return SecCheck{Result: "可登录的默认/系统账号：" + strings.Join(login, "、"), Judge: "不符合"}
		}
		return SecCheck{Result: "默认账号存在但不可登录（shell 为 nologin/false）：" + strings.Join(have, "、"), Judge: "符合"}
	})

	//4. 口令复杂度 —— 用已采集的 pwquality 8 项 + PAM 内容
	add("口令复杂度策略", func() SecCheck {
		q := d.Quality
		vals := []string{q.Minlen, q.Dcredit, q.Ucredit, q.Lcredit, q.Ocredit, q.Minclass, q.Maxrepeat, q.Maxsequence}
		set := 0
		var show []string
		names := []string{"minlen", "dcredit", "ucredit", "lcredit", "ocredit", "minclass", "maxrepeat", "maxsequence"}
		for i, v := range vals {
			if strings.TrimSpace(v) != "" {
				set++
				show = append(show, names[i]+"="+strings.TrimSpace(v))
			}
		}
		pamUse := strings.Contains(d.PamSystem, "pam_pwquality") || strings.Contains(d.PamSystem, "pam_cracklib")
		res := "/etc/security/pwquality.conf 生效项：" + strOr(strings.Join(show, " "), "无（全部注释，按默认值）")
		if pamUse {
			res += "\n/etc/pam.d/system-auth 已引用 pam_pwquality/pam_cracklib"
		} else {
			res += "\n/etc/pam.d/system-auth 未引用 pam_pwquality/pam_cracklib"
		}
		switch {
		case pamUse && set > 0:
			return SecCheck{Result: res, Judge: "符合"}
		case pamUse:
			return SecCheck{Result: res, Judge: "部分符合"}
		}
		return SecCheck{Result: res, Judge: "不符合"}
	})

	//5. 登录失败处理 —— 用已采集的 /etc/pam.d/system-auth、/etc/pam.d/sshd 内容
	add("登录失败处理", func() SecCheck {
		src := d.PamSystem + "\n" + d.PamSSH
		hitPam := hitAny(src, "pam_tally2", "pam_faillock")
		hitDeny := hitAny(src, "deny=")
		if !hitPam {
			return SecCheck{Result: "已采集的 PAM 配置中无 pam_tally2 / pam_faillock，登录失败无锁定策略", Judge: "不符合"}
		}
		var lines []string
		for _, l := range strings.Split(src, "\n") {
			if hitAny(l, "pam_tally2", "pam_faillock", "deny=", "unlock_time") {
				lines = append(lines, strings.TrimSpace(l))
			}
		}
		res := trunc(strings.Join(lines, "\n"), 600)
		if hitDeny {
			return SecCheck{Result: res, Judge: "符合"}
		}
		return SecCheck{Result: res + "\n（未见 deny= 阈值，锁定次数可能为默认）", Judge: "部分符合"}
	})

	//6. 会话超时 —— 用已采集的 /etc/profile 环境变量内容
	add("会话超时 TMOUT", func() SecCheck {
		if strings.TrimSpace(d.Env) == "" {
			return SecCheck{Result: "/etc/profile 未采集到内容", Judge: "采集失败"}
		}
		var lines []string
		for _, l := range strings.Split(d.Env, "\n") {
			if strings.Contains(l, "TMOUT") {
				lines = append(lines, strings.TrimSpace(l))
			}
		}
		if len(lines) == 0 {
			return SecCheck{Result: "/etc/profile 有效配置中无 TMOUT，会话不会自动超时退出", Judge: "不符合"}
		}
		return SecCheck{Result: strings.Join(lines, "\n"), Judge: "符合"}
	})

	//7. 日志外发 —— 用已采集的 /etc/rsyslog.conf
	add("日志外发", func() SecCheck {
		if strings.TrimSpace(d.Rsyslog) == "" {
			return SecCheck{Result: "/etc/rsyslog.conf 未采集到内容", Judge: "采集失败"}
		}
		var fwd []string
		for _, l := range strings.Split(d.Rsyslog, "\n") {
			if strings.Contains(l, "@@") || strings.HasPrefix(strings.TrimSpace(l), "@") {
				fwd = append(fwd, strings.TrimSpace(l))
			}
		}
		if len(fwd) == 0 {
			return SecCheck{Result: "/etc/rsyslog.conf 无 @@ 转发，日志仅落本机", Judge: "部分符合"}
		}
		return SecCheck{Result: strings.Join(fwd, "\n"), Judge: "符合"}
	})

	//8. 日志保留 —— 用已采集的日志切割配置
	add("日志保留与切割", func() SecCheck {
		if strings.TrimSpace(d.Logrotate) == "" {
			return SecCheck{Result: "日志切割配置未采集到", Judge: "采集失败"}
		}
		var lines []string
		for _, l := range strings.Split(d.Logrotate, "\n") {
			if hitAny(l, "rotate", "weekly", "daily", "monthly") {
				lines = append(lines, strings.TrimSpace(l))
			}
		}
		if len(lines) == 0 {
			return SecCheck{Result: "未找到 rotate/weekly/daily 等保留策略", Judge: "不符合"}
		}
		return SecCheck{Result: trunc(strings.Join(lines, "\n"), 500), Judge: "需人工确认"}
	})

	//9. 审计规则 —— 用已采集的 auditctl -l 结果
	add("审计规则", func() SecCheck {
		o := strings.TrimSpace(d.AuditCtl)
		if o == "" || strings.EqualFold(o, "No rules") {
			return SecCheck{Result: "auditctl -l 无任何规则（审计服务未启用或未配置规则）", Judge: "不符合"}
		}
		n := 0
		for _, l := range strings.Split(o, "\n") {
			if strings.TrimSpace(l) != "" {
				n++
			}
		}
		return SecCheck{Result: fmt.Sprintf("共 %d 条规则：\n%s", n, trunc(o, 600)), Judge: "符合"}
	})

	//10. 高危服务端口 —— 用已采集的端口列表
	add("高危服务端口", func() SecCheck {
		if len(d.Port) == 0 {
			return SecCheck{Result: "端口列表未采集到", Judge: "采集失败"}
		}
		var hit []string
		for _, p := range d.Port {
			local := strings.TrimSpace(p.Local)
			for _, port := range []string{":21", ":23", ":69", ":512", ":513", ":514"} {
				if strings.HasSuffix(local, port) {
					hit = append(hit, fmt.Sprintf("%s %s %s", p.Netid, local, p.Process))
					break
				}
			}
		}
		if len(hit) == 0 {
			return SecCheck{Result: fmt.Sprintf("端口列表共 %d 项，未发现 telnet/ftp/tftp/rsh 端口", len(d.Port)), Judge: "符合"}
		}
		return SecCheck{Result: "发现高危端口：\n" + strings.Join(hit, "\n"), Judge: "不符合"}
	})

	//11. 防病毒 / 主机防护 —— 用已采集的进程列表
	add("防病毒与主机防护", func() SecCheck {
		if strings.TrimSpace(d.PS) == "" {
			return SecCheck{Result: "进程列表未采集到", Judge: "采集失败"}
		}
		var hit []string
		for _, l := range strings.Split(d.PS, "\n") {
			if hitAny(l, "clamd", "clamav", "safedog", "aegis", "edr", "qax", "dhagent", "sangfor", "ossec", "wazuh", "tianqing") {
				hit = append(hit, trunc(strings.TrimSpace(l), 160))
			}
		}
		if len(hit) == 0 {
			return SecCheck{Result: "进程列表中未发现防病毒/主机防护类软件（命令行判定，需结合主机侧确认）", Judge: "不符合"}
		}
		return SecCheck{Result: strings.Join(hit, "\n"), Judge: "符合"}
	})

	//12. sudo 授权 —— 用已采集的 /etc/sudoers
	add("sudo 授权", func() SecCheck {
		var lines []string
		for _, l := range strings.Split(d.Sudoers, "\n") {
			t := strings.TrimSpace(l)
			if t == "" || strings.HasPrefix(t, "Defaults") || strings.HasPrefix(t, "File:") {
				continue
			}
			lines = append(lines, t)
		}
		if len(lines) == 0 {
			return SecCheck{Result: "/etc/sudoers 无实际授权项（仅有 Defaults 等全局配置）", Judge: "不适用"}
		}
		return SecCheck{Result: trunc(strings.Join(lines, "\n"), 800), Judge: "需人工确认"}
	})

	//13. 内核与补丁 —— 用已采集的版本信息
	add("内核与系统版本", func() SecCheck {
		v := strings.TrimSpace(d.Info.Version)
		if v == "" {
			return SecCheck{Result: "系统版本未采集到", Judge: "采集失败"}
		}
		return SecCheck{Result: trunc(v, 400), Judge: "需人工确认"}
	})

	//14. SSH 生效算法 —— 需要 sshd -T（现有采集只有 sshd_config 的 6 个字段）
	add("SSH 生效配置与算法", func() SecCheck {
		if c == nil {
			return SecCheck{Result: "本机模式未采集 sshd -T", Judge: "采集失败"}
		}
		o := strings.TrimSpace(runCmd(`sshd -T 2>/dev/null | grep -Ei "^(permitrootlogin|passwordauthentication|permitemptypasswords|maxauthtries|ciphers|macs|kexalgorithms)" | head -n 10`, c))
		if o == "" {
			return SecCheck{Result: "sshd -T 无输出（未安装 openssh-server 或权限不足）", Judge: "需人工确认"}
		}
		weak := []string{"3des-cbc", "aes128-cbc", "aes192-cbc", "aes256-cbc", "arcfour",
			"hmac-md5", "diffie-hellman-group1-sha1", "diffie-hellman-group14-sha1",
			"diffie-hellman-group-exchange-sha1", "ssh-dss"}
		found := make([]string, 0, 4)
		for _, w := range weak {
			if strings.Contains(o, w) {
				found = append(found, w)
			}
		}
		if len(found) > 0 {
			return SecCheck{Result: trunc(o, 700) + "\n检测到弱算法：" + strings.Join(found, "、"), Judge: "部分符合"}
		}
		return SecCheck{Result: trunc(o, 700), Judge: "符合"}
	})

	//15. 时间同步 —— 需要额外命令
	add("时间同步", func() SecCheck {
		if c == nil {
			return SecCheck{Result: "本机模式未采集时间同步", Judge: "采集失败"}
		}
		o := strings.TrimSpace(runCmd(`timedatectl 2>/dev/null | grep -Ei "synchronized|NTP service|System clock"; chronyc tracking 2>/dev/null | head -n 3; ntpq -p 2>/dev/null | head -n 3`, c))
		if o == "" {
			return SecCheck{Result: "未采集到时间同步信息（无 timedatectl / chrony / ntp）", Judge: "需人工确认"}
		}
		if hitAny(o, "synchronized: yes", "Leap status     : Normal", "Leap status : Normal") {
			return SecCheck{Result: o, Judge: "符合"}
		}
		return SecCheck{Result: o, Judge: "部分符合"}
	})

	return out
}

func strOr(s string, def string) string {
	if strings.TrimSpace(s) == "" {
		return def
	}
	return s
}

// hitAny 判断 s 是否包含任意一个子串
func hitAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if sub != "" && strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// trunc 安全截断（长度不足时原样返回，不会越界）
func trunc(s string, n int) string {
	if n <= 0 || len(s) <= n {
		return s
	}
	return s[:n] + "\n...（已截断）"
}
