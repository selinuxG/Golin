package network

import (
	"bufio"
	"fmt"
	"golin/global"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Rsyslog(protocol, port string) {
	intport, err := strconv.Atoi(port)
	if err != nil {
		zlog.Warn("syslog端口转换失败！")
		return
	}

	if intport > 65535 || intport == 0 {
		zlog.Warn("非法端口！")
		return
	}
	if protocol == "tcp" || protocol == "udp" {
		Rsysrun(protocol, port)
		return

	}
	zlog.Warn("协议只能是tcp/udp！")
}

func Rsysrun(protocol, port string) {
	fmt.Printf("Listening on %s:%s\n", protocol, port)

	// 监听地址
	address := fmt.Sprintf(":%s", port)
	// 创建监听器
	var listener net.Listener
	var packetConn net.PacketConn
	var err error
	if protocol == "tcp" {
		listener, err = net.Listen("tcp", address)
		if err != nil {
			panic(err)
		}
		defer listener.Close()
		for {
			// 接收连接
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("Error accepting connection: %s\n", err.Error())
				continue
			}

			// 持续接收客户端发送的消息
			for {
				buffer := make([]byte, 1024)
				n, err := conn.Read(buffer)
				if err != nil {
					fmt.Printf("Error reading from connection: %s\n", err.Error())
					break
				}
				name := fmt.Sprintf("tcp_%s", conn.RemoteAddr().String())
				err = AppendToFile(name, string(buffer[:n]))
				if err != nil {
					zlog.Warn("rsyslog写入日志失败")
				}
				fmt.Printf("接受消息来源 %s:\n消息内容:%s\n", conn.RemoteAddr().String(), string(buffer[:n]))
			}

		}
	} else {
		packetConn, err = net.ListenPacket("udp", address)
		if err != nil {
			panic(err)
		}
		defer packetConn.Close()
		for {
			// 接收数据包
			buffer := make([]byte, 1024)
			n, addr, err := packetConn.ReadFrom(buffer)
			if err != nil {
				fmt.Printf("Error reading from connection: %s\n", err.Error())
				continue
			}
			name := fmt.Sprintf("udp_%s", addr.String())
			err = AppendToFile(name, string(buffer[:n]))
			if err != nil {
				zlog.Warn("rsyslog写入日志失败")
			}
			fmt.Printf("接受消息来源 %s\n消息内容:%s\n", addr.String(), string(buffer[:n]))
		}
	}
}

func AppendToFile(filename string, content string) error {
	filename = strings.Split(filename, ":")[0]
	filename += ".log"
	//先判断目录是否存在
	appendfile := filepath.Join(global.Succpath, "rsyslog")
	_, err := os.Stat(appendfile)
	if err != nil {
		os.MkdirAll(appendfile, os.FileMode(global.FilePer))
	}
	// 检测文件是否存在
	createfile := filepath.Join(appendfile, filename)
	//追加写入文件
	_, err = os.Stat(createfile)
	if os.IsNotExist(err) {
		os.Create(createfile)
	}
	file, err := os.OpenFile(createfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(global.FilePer))
	defer file.Close()
	if err != nil {
		return err
	}
	write := bufio.NewWriter(file)
	write.WriteString(content + "\n")
	err = write.Flush()
	if err != nil {
		return err
	}
	return nil
}
