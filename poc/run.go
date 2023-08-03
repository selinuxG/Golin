package poc

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

type Flagcve struct {
	url  string
	cve  string
	flag string
}

func CheckPoc(url, app string) {

	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}

	app = strings.ToLower(app)
	switch {
	case strings.Contains(app, "spring"):
		CVE_2022_22947(url, "pwd")
	}

}

func echoFlag(flag Flagcve) {
	fmt.Printf("\033[2K\r") // 擦除整行
	fmt.Printf("\r| %-2s | %-15s | %-15s |%s\n",
		fmt.Sprintf("%s", color.RedString("%s", "√")),
		fmt.Sprintf("%s", color.RedString("发现漏洞_%s", flag.cve)),
		fmt.Sprintf("%s", color.RedString(flag.url)),
		fmt.Sprintf("%s", color.RedString(flag.flag)),
	)

}
