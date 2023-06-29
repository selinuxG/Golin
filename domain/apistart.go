package domain

import (
	"fmt"
)

func apiStart(domain string) {

	_, _ = fmt.Fprintf(colorOutput, "\n%s[*] 开始调用FOFA 最大获取条数：%d %s  \n", greenColor, size, resetColor)
	fofaApi(domain)

	_, _ = fmt.Fprintf(colorOutput, "\n%s[*] 开始调用RapidDNS 最大获取条数：%d %s\n", greenColor, size, resetColor)
	rapidDNS(domain)

	fmt.Printf("\n")
}
