// link:d-eyes

package clientinfo

import (
	"bufio"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func Netstat(cmd *cobra.Command, args []string) {
	networkData := make([][]string, 0)
	var remoteIp []string
	ps, err := process.Processes()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, p := range ps {

		pid := os.Getpid()
		if pid == int(p.Pid) || p.Pid == 0 {
			continue
		}

		connList := make([]string, 0)
		connection := make([]string, 0)
		_pc, _ := p.Connections()
		for _, conn := range _pc {
			if conn.Family == 1 {
				continue
			}
			c := fmt.Sprintf(
				"%v:%v<->%v:%v(%v)\n",
				conn.Laddr.IP, conn.Laddr.Port, conn.Raddr.IP, conn.Raddr.Port, conn.Status,
			)
			remoteIp = append(remoteIp, conn.Raddr.IP)

			connection = append(connection, c)
		}
		_pUname, _ := p.Username()
		if len(connection) > 0 && _pUname != "" {
			network := strings.Join(connection, "")
			_exe, _ := p.Exe()
			path := StringNewLine(_exe, 25)
			username, _ := p.Username()
			connList = append(connList, fmt.Sprintf("%v", p.Pid), fmt.Sprintf("%v", username), network, path)
			networkData = append(networkData, connList)
		}
	}

	tableConn := tablewriter.NewWriter(os.Stdout)
	tableConn.SetHeader([]string{"pid", "user", "local/remote(TCP Status)", "program name"})
	tableConn.SetBorder(true)
	tableConn.SetRowLine(true)
	tableConn.AppendBulk(networkData)
	tableConn.Render()
	remoteIpNew := RemoveRepeatedElement(remoteIp)

	if len(remoteIpNew) > 0 {
		if err = WriteSliceToFile(remoteIpNew, "netstat.txt"); err == nil {
			fmt.Println("[*] 建立连接IP以写入到netstat.txt文件中...")
		}
	}
}

func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		if arr[i] == "127.0.0.1" || arr[i] == "0.0.0.0" || arr[i] == "::" || arr[i] == "::1" || arr[i] == "" {
			continue
		}
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

func StringNewLine(str string, ln uint8) string {
	var sub_str string
	res_str := ""
	for {
		if len(str) < int(ln) {
			res_str += str
			break
		}
		sub_str = str[0:ln]
		str = str[ln:]
		res_str += sub_str + "\n"
	}
	return res_str
}

// WriteSliceToFile 将切片中的每个元素写入指定的文件
func WriteSliceToFile(slice []string, filename string) error {
	// 打开文件，如果文件不存在则创建，如果文件存在则清空并打开。
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	// 使用 bufio.NewWriter 创建一个新的 Writer 对象。
	writer := bufio.NewWriter(file)

	// 遍历切片，将每个元素写入文件。
	for _, item := range slice {
		_, err := writer.WriteString(item + "\n") // 在每个元素后添加换行符。
		if err != nil {
			return fmt.Errorf("写入文件失败: %v", err)
		}
	}

	// 确保所有缓冲的数据都已写入底层 io.Writer。
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("刷新到文件失败: %v", err)
	}

	return nil
}
