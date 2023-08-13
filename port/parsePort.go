package port

import (
	"strconv"
	"strings"
)

func parsePort(port string) {
	if len(portlist) == 1 { //如果是快速扫描则已经有端口了
		return
	}

	if port == "0" {
		portlist = default_port
		return
	}

	if strings.Count(port, ",") == 0 {
		port = port + ","
	}

	for _, p := range strings.Split(port, ",") {
		if p == "" {
			continue
		}

		if strings.Count(p, "-") == 1 { //范围
			start := strings.Split(p, "-")[0]
			end := strings.Split(p, "-")[1]
			startNum, _ := strconv.Atoi(start)
			endNum, _ := strconv.Atoi(end)
			if startNum > endNum {
				continue
			}
			if startNum == endNum {
				portlist = append(portlist, start)
				continue
			}
			for i := startNum; i <= endNum; i++ {
				portlist = append(portlist, strconv.Itoa(i))
			}
		}

		if strings.Count(p, "-") == 0 {
			pNun, _ := strconv.Atoi(p)
			if pNun < 65535 {
				portlist = append(portlist, p)
			}
		}

	}
}

var (
	default_port = []string{
		"22",  // SSH
		"80",  // HTTP
		"81",  // HTTP
		"82",  // HTTP
		"88",  // HTTP
		"443", // HTTPS
		"21",  // FTP
		"8888",
		"9999",
		"8250",
		"8650",
		"8750",
		"9200",
		"9080",
		"4430",
		"8989",
		"808",
		"8443",
		"9080",
		"8001",
		"8034",
		"10002",
		"25",    // SMTP
		"110",   // POP3
		"143",   // IMAP
		"465",   //SMTPS
		"587",   //Submission
		"993",   //IMAPS
		"995",   //POP3S
		"1080",  //SOCKS
		"1194",  //开放VPN
		"5900",  //VNC
		"2181",  //ZooKeeper
		"53",    // DNS
		"3389",  // RDP
		"23",    // Telnet
		"514",   //Syslog
		"389",   // LDAP
		"1433",  // MSSQL
		"1521",  // Oracle Database
		"6379",  // Redis
		"27017", // MongoDB
		"8000",  //web
		"8080",  // HTTP (Alternate)
		"8443",  // HTTPS (Alternate)
		"9000",  // PHP-FPM
		"161",   // SNMP Trap
		"162",   // SNMP Trap
		"445",   // SMB
		"137",   // NetBIOS
		"138",   // NetBIOS
		"139",   // NetBIOS
		"1434",  // MSSQL (Alternate)
		"1723",  //PPTP
		"3306",  // MySQL (Alternate)
		"5432",  // PostgreSQL (Alternate)
		"6378",  // redis
		"5601",  //kibana
		"1080",  //sock
		"1194",  //vpn
		"5900",  //vnc
		"5901",  //vnc
		"6066",
		"8085",
		"7105",
		"5672", //RabbitMq
		"6000", //x11
		"6443", //K8S
		"9000", //Hadoop
		"3000",
		"8001",
		"8002",
		"8003",
		"8004",
		"8005",
		"8006",
		"5984",
		"8007",
		"8008",
		"8009",
		"8010",
		"8081",
		"8082",
		"8083",
		"8084",
		"8085",
		"8070",
		"7443",
		"8161",
		"9999",
		"8899",
		"8989",
		"9999",
		"9001",
		"9002",
		"9003",
		"9010",
		"9090",
		"9099",
		"10808",
		"10809",
		"61616",
	}
)
