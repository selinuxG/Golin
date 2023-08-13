package poc

import (
	"fmt"
	"net"
	"time"
)

func ZookeeperCon(host, port string) {
	address := fmt.Sprintf("%s:%s", host, port)
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte("envi"))
	if err != nil {
		return
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return
	}

	response := string(buffer[:n])
	// 如果响应中包含 "Environment"，则是 Zookeeper 服务
	if n > 0 && response[0:11] == "Environment" {
		flags := Flagcve{fmt.Sprintf("%s:%s", host, port), "Zookeeper未授权访问", "可通过nc参数执行conf、envi等命令验证"}
		echoFlag(flags)
	}
}
