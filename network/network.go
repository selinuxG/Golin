package network

import (
	"bytes"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"github.com/spf13/cobra"
	"golin/config"
	"log"
	"os"
	"strings"
	"time"
)

var (
	zlog          = config.Log
	echointerface string
	interfaceall  bool
	bfp           string
	err           error
	ifaceslist    []string
)

func Networkrun(cmd *cobra.Command, args []string) {
	syslog, err := cmd.Flags().GetBool("syslog")
	if syslog {
		if len(args) == 2 {
			Rsyslog(args[0], args[1])
			return
		}
		zlog.Warn("参数错误,接受2个参数,第一个是协议（tcp/udp）,第二个是端口")
		return
	}

	interfaceall, err = cmd.Flags().GetBool("interfaceall")
	if err != nil {
		zlog.Warn("解析interfaceall参数失败,退出！")
		return
	}
	// 查找所有可用的网络接口
	if interfaceall {
		Interfaceall()
		return
	}
	//解析读取网卡
	echointerface, err = cmd.Flags().GetString("interface")
	if err != nil {
		zlog.Warn("解析address参数失败,退出！")
		return
	}
	ifaces, err := pcap.FindAllDevs()
	if err != nil {
		panic(err)
	}
	// 遍历所有接口并输追加到切片中去
	for _, iface := range ifaces {
		ifaceslist = append(ifaceslist, iface.Name)
	}
	if !InSlice(ifaceslist, echointerface) {
		zlog.Warn("网卡不存在，通过 -a 可以获取所有网卡,退出！")
		return
	}
	//获取监听语句
	bfp, err = cmd.Flags().GetString("bfp")
	if err != nil {
		zlog.Warn("解析bfp参数失败,退出！")
		return
	}
	if bfp == "" {
		zlog.Warn("bfp参数为空，我也不知道你要监听什么,退出！")
		return
	}
	run(echointerface, bfp)
}

func run(i string, bfp string) {
	fmt.Printf("------------------------正在监听网卡:%s 监听类型:%s\n", i, bfp)
	if handle, err := pcap.OpenLive(i, 65536, true, -1*time.Second); err != nil {
		log.Fatal(err)
	} else {
		defer handle.Close()
		if err := handle.SetBPFFilter(bfp); err != nil {
			log.Fatal(err)
		}
		// create new file for writing captured packets
		file, err := os.Create("network.pcap")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		// create new pcap writer with file as output
		writer := pcapgo.NewWriter(file)
		writer.WriteFileHeader(65536, layers.LinkTypeEthernet)
		// capture packets and write to file
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			//fmt.Printf("\n\n--------------------\n")
			if err := writer.WritePacket(packet.Metadata().CaptureInfo, packet.Data()); err != nil {
				log.Fatal(err)
			}
			// 传输层协议判断
			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				fmt.Println("\n--------------------传输层协议类型：TCP")
			} else if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
				fmt.Println("\n--------------------传输层协议类型：UDP")
			}
			//应用层协议判断
			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)
				srcPort := tcp.SrcPort
				dstPort := tcp.DstPort
				payload := tcp.Payload
				if len(payload) > 0 {
					if isSSH(payload) {
						fmt.Println("--------------------应用层协议类型：SSH")
					} else if isMySQL(payload) {
						fmt.Println("--------------------应用层协议类型：MySQL")
						applicationLayer := packet.ApplicationLayer()
						if applicationLayer != nil {
							payload := applicationLayer.Payload()
							fmt.Printf("Payload: %s\n", string(payload))
						}
					} else if isHTTP(payload) {
						fmt.Println("--------------------应用层协议类型：HTTP")
					} else if srcPort == 443 || dstPort == 443 {
						fmt.Println("--------------------应用层协议类型：HTTPS")
					} else if isFTP(payload) {
						fmt.Println("--------------------应用层协议类型：FTP")
					} else if isSMB(payload) {
						fmt.Println("--------------------应用层协议类型：SMB")
					} else if isDHCP(payload) {
						fmt.Println("--------------------应用层协议类型：DHCP")
					} else if isPing(payload) {
						fmt.Println("--------------------应用层协议类型：ICMP")
					}
				}
				// 获取负载payload
				//applicationLayer := packet.ApplicationLayer()
				//if applicationLayer != nil {
				//	payload := applicationLayer.Payload()
				//	fmt.Printf("Payload: %s\n", string(payload))
				//}
				// 获取网络层和传输层信息
				networkLayer := packet.NetworkLayer()
				transLayer := packet.TransportLayer()
				if networkLayer != nil && transLayer != nil {
					// 输出IP地址
					if ipLayer, ok := networkLayer.(*layers.IPv4); ok {
						fmt.Println("来源IP地址:", ipLayer.SrcIP, "目的IP地址:", ipLayer.DstIP)
					}
					// 输出端口信息
					if tcpLayer, ok := transLayer.(*layers.TCP); ok {
						fmt.Println("源端口:", tcpLayer.SrcPort, "目标端口:", tcpLayer.DstPort)
					} else if udpLayer, ok := transLayer.(*layers.UDP); ok {
						fmt.Println("源端口:", udpLayer.SrcPort, "目标端口:", udpLayer.DstPort)
					}
				}
			}
		}
	}
}

func Interfaceall() {
	ifaces, err := pcap.FindAllDevs()
	if err != nil {
		panic(err)
	}
	// 遍历所有接口并输出它们的名称
	for _, iface := range ifaces {
		fmt.Println("网卡名称：", iface.Name)
		for _, address := range iface.Addresses {
			fmt.Println(address.IP)
		}
		fmt.Println("----------------------")
	}
}

// InSlice 判断字符串是否在 slice 中。
func InSlice(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// 判断负载中是否包含SSH协议标识符
func isSSH(payload []byte) bool {
	payloadStr := strings.ToLower(string(payload))
	sqlKeywords := []string{"ssh", "SSH"}
	for _, keyword := range sqlKeywords {
		if strings.Contains(payloadStr, keyword) {
			return true
		}
	}
	return false
}

// 判断负载中是否包含MySQL协议标识符
func isMySQL(payload []byte) bool {
	// 判断负载中是否包含特定的命令或响应码
	payloadStr := strings.ToLower(string(payload))
	sqlKeywords := []string{"select", "insert", "update", "delete", "create", "drop", "alter", "show", "variabl", "version()"}
	for _, keyword := range sqlKeywords {
		if strings.Contains(payloadStr, keyword) {
			return true
		}
	}
	if bytes.HasPrefix(payload, []byte{0x00, 0x01, 0x00, 0x00, 0x00}) {
		return true
	}
	return false
}

// 判断负载中是否包含HTTP协议标识符
func isHTTP(payload []byte) bool {
	// 判断负载中是否包含HTTP请求或响应的头部信息等等
	return bytes.Contains(payload, []byte("\r\nHost:")) ||
		bytes.Contains(payload, []byte("\r\nUser-Agent:")) ||
		bytes.Contains(payload, []byte("\r\nAccept:")) ||
		bytes.Contains(payload, []byte("\r\nContent-Length:")) ||
		bytes.Contains(payload, []byte("HTTP/1."))
}

// 判断是否为FTP协议
func isFTP(payload []byte) bool {
	if len(payload) < 4 {
		return false
	}
	return bytes.Equal(payload[:4], []byte("USER")) || bytes.Equal(payload[:4], []byte("PASS"))
}

// 判断是否为SMB协议
func isSMB(payload []byte) bool {
	if len(payload) < 4 {
		return false
	}
	return bytes.Equal(payload[:4], []byte{0xff, 0x53, 0x4d, 0x42})
}

// 判断是否为DHCP协议
func isDHCP(payload []byte) bool {
	if len(payload) < 240 {
		return false
	}
	return bytes.Equal(payload[236:240], []byte{0x63, 0x82, 0x53, 0x63})
}

// 判断是否为ICMP协议
func isPing(payload []byte) bool {
	// 判断 payload 长度是否大于 28（ICMP ping 包最少 28 个字节）
	if len(payload) < 28 {
		return false
	}

	// 判断是否为 ICMP ping 包
	if payload[0] == 0x08 && payload[1] == 0x00 {
		return true
	}

	return false
}
