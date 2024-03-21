package scan

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"golin/poc"
	"golin/scan/crack"
	"os"
	"strconv"
	"strings"
)

// echoCrack 输出弱口令的资产信息
func echoCrack() {
	if len(crack.MapCrackHost) <= 0 {
		return
	}

	var data [][]string

	for _, sussCrack := range crack.MapCrackHost {
		data = append(data, []string{sussCrack.Host, strconv.Itoa(sussCrack.Port), sussCrack.User, sussCrack.Passwd, sussCrack.Mode})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Host", "Port", "User", "Passwd", "Mode"})
	table.AppendBulk(data)
	fmt.Printf(clearLine)
	table.Render()
}

// echoPoc 输出漏洞资产信息
func echoPoc() {
	if len(poc.ListPocInfo) <= 0 {
		return
	}
	var data [][]string
	for _, sussPoc := range poc.ListPocInfo {
		sussPoc.Cve = strings.ReplaceAll(sussPoc.Cve, "poc-yaml-", "")
		data = append(data, []string{sussPoc.Url, sussPoc.Cve, sussPoc.Flag})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Url", "Name", "Info"})
	table.AppendBulk(data)
	fmt.Printf(clearLine)
	table.Render()

}
