package domain

import (
	"fmt"
	"github.com/fatih/color"
)

func apiStart(domain string) {

	fmt.Printf("\n[*] 开始调用FOFA 获取条数:%s\n",
		color.GreenString("%d", size),
	)
	fofaApi(domain)

	fmt.Printf("\n[*] 开始调用RapidDNS 获取条数:%s\n",
		color.GreenString("%d", size),
	)
	rapidDNS(domain)

	fmt.Printf("\n")
}
