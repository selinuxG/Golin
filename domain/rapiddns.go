package domain

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/parnurzeal/gorequest"
	"os"
	"strconv"
	"time"
)

type rapid struct {
	Data []struct {
		Date      string    `json:"date"`
		Name      string    `json:"name"`
		Timestamp time.Time `json:"timestamp"`
		Type      string    `json:"type"`
		Value     string    `json:"value"`
	} `json:"data"`
	Description string `json:"description"`
	Maxpage     int    `json:"maxpage"`
	Status      string `json:"status"`
	Total       int    `json:"total"`
}

func rapidDNS(domain string) {

	request := gorequest.New().Timeout(10 * time.Second)
	if request == nil {
		fmt.Printf("[-] 调用RapidDNS失败！\n")
		return
	}
	domain = fmt.Sprintf("https://rapiddns.io/api/v1/%s?size=%d&page=1", domain, size)
	req := request.Get(domain)
	if req == nil {
		fmt.Printf("[-] 调用RapidDNS失败！\n")
		return
	}

	resp, body, errs := req.End()
	if len(errs) > 0 {
		fmt.Printf("[-] 调用RapidDNS失败！\n")
		return
	}
	defer resp.Body.Close()

	domainList := &rapid{}
	err := json.Unmarshal([]byte(body), domainList)
	if err != nil {
		fmt.Printf("[-] 调用RapidDNS失败！\n")
		return
	}
	if domainList.Status != "200" {
		fmt.Printf("[-] 调用RapidDNS失败！\n")
		return
	}

	if len(domainList.Data) == 0 {
		fmt.Printf("[-] RapidDNS 未收录此数据...\n")
		return
	}

	var data [][]string
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"doamin", "ip", "type", "date"})

	for _, da := range domainList.Data {
		row := []string{da.Name, da.Value, da.Type, da.Date}
		data = append(data, row)
	}

	table.SetFooter([]string{"", "", "Total", strconv.Itoa(len(domainList.Data))}) // Add Footer
	table.AppendBulk(data)                                                         // Add Bulk Data
	table.Render()

}
