package domain

import (
	"fmt"
)

func apiStart(domain string) {

	_, _ = fmt.Fprintf(colorOutput, "\n%s[*] 开始调用FOFA %s\n", greenColor, resetColor)
	fofaApi(domain)

	_, _ = fmt.Fprintf(colorOutput, "\n%s[*] 开始调用RapidDNS %s\n", greenColor, resetColor)
	rapidDNS(domain)

	fmt.Printf("\n")
}
