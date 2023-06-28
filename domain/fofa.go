package domain

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/parnurzeal/gorequest"
	"os"
	"time"
)

type FoFA struct {
	Error           bool       `json:"error"`
	ConsumedFpoint  int        `json:"consumed_fpoint"`
	RequiredFpoints int        `json:"required_fpoints"`
	Size            int        `json:"size"`
	Page            int        `json:"page"`
	Mode            string     `json:"mode"`
	Query           string     `json:"query"`
	Results         [][]string `json:"results"`
}

func fofa_Api(domain string) {
	request := gorequest.New().Timeout(10 * time.Second)
	if request == nil {
		echoerr()
		return
	}
	req := request.Get(domain)
	if req == nil {
		echoerr()
		return
	}

	resp, body, errs := req.End()
	if len(errs) > 0 {
		echoerr()
		return
	}
	defer resp.Body.Close()

	domainList := &FoFA{}
	err := json.Unmarshal([]byte(body), domainList)
	if err != nil {
		echoerr()
		return
	}
	if domainList.Error {
		echoerr()
		return
	}
	if len(domainList.Results) > 0 {
		var data [][]string
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"doamin", "ip", "port"})

		for _, result := range domainList.Results {
			row := []string{result[0], result[1], result[2]}
			data = append(data, row)
		}
		for _, v := range data {
			table.Append(v)
		}
		table.Render()
	}
}

func echoerr() {
	fmt.Println("[-] 调用FOFA-API失败！")
}
