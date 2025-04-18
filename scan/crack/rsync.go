package crack

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func Rsync(host, port string) {
	timeout := 3 * time.Second

	// 建立连接
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", host, port), 3*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1024)

	// 1. 读取服务器初始greeting
	n, err := conn.Read(buffer)
	if err != nil {
		return
	}

	greeting := string(buffer[:n])

	// 获取服务器版本号
	version := strings.TrimSpace(strings.TrimPrefix(greeting, "@RSYNCD:"))

	// 2. 回应相同的版本号
	_, err = conn.Write([]byte(fmt.Sprintf("@RSYNCD: %s\n", version)))
	if err != nil {
		return
	}

	// 3. 选择模块 - 先列出可用模块
	_, err = conn.Write([]byte("#list\n"))
	if err != nil {
		return
	}

	// 4. 读取模块列表
	var moduleList strings.Builder
	for {
		n, err = conn.Read(buffer)
		if err != nil {
			break
		}
		chunk := string(buffer[:n])
		moduleList.WriteString(chunk)
		if strings.Contains(chunk, "@RSYNCD: EXIT") {
			break
		}
		fmt.Println("卡住了")
	}

	modules := strings.Split(moduleList.String(), "\n")
	for _, module := range modules {
		if strings.HasPrefix(module, "@RSYNCD") || module == "" {
			continue
		}

		// 获取模块名
		moduleName := strings.Fields(module)[0]

		// 5. 为每个模块创建新连接尝试认证
		authConn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", host, port), timeout)
		if err != nil {
			continue
		}

		// 重复初始握手
		_, err = authConn.Read(buffer)
		if err != nil {
			_ = authConn.Close()
			continue
		}

		_, err = authConn.Write([]byte(fmt.Sprintf("@RSYNCD: %s\n", version)))
		if err != nil {
			_ = authConn.Close()
			continue
		}

		// 6. 选择模块
		_, err = authConn.Write([]byte(moduleName + "\n"))
		if err != nil {
			_ = authConn.Close()
			continue
		}

		// 7. 等待认证挑战
		n, err = authConn.Read(buffer)

		if err != nil {
			_ = authConn.Close()
			continue
		}

		authResponse := string(buffer[:n])
		if strings.Contains(authResponse, "@RSYNCD: OK") {
			_ = authConn.Close()
			intport, _ := strconv.Atoi(port)
			end(host, "未授权访问", "", intport, "Rsync")
		}
	}
}
