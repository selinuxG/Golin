package port

import (
	"strconv"
	"strings"
)

func parsePort(port string) {
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
		"20",    //ftp
		"22",    // SSH
		"80",    // HTTP
		"443",   // HTTPS
		"21",    // FTP
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
		"53",    // DNS
		"3389",  // RDP
		"23",    // Telnet
		"514",   //Syslog
		"389",   // LDAP
		"1433",  // MSSQL
		"1521",  // Oracle Database
		"6379",  // Redis
		"27017", // MongoDB
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
		"3306",  // MySQL (Alternate)
		"5432",  // PostgreSQL (Alternate)
		"6378",  // redis
		"9200",  //es
		"9300",  //es
		"5601",  //kibana
		"1080",  //sock
		"1194",  //vpn
		"5900",  //vnc
	}
)
