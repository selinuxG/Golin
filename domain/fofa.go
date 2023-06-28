package domain

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/parnurzeal/gorequest"
	"os"
	"strconv"
	"strings"
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
	base64query := fmt.Sprintf("%s && is_domain=true", domain)
	size := os.Getenv("FOFA_SIZE")
	if size == "" {
		size = strconv.Itoa(100)
	}
	domain = fmt.Sprintf("https://fofa.info/api/v1/search/all?email=%s&key=%s&size=%s&qbase64=%s", os.Getenv("FOFA_USER"), os.Getenv("FOFA_KEY"), size, base64.StdEncoding.EncodeToString([]byte(base64query)))
	//fmt.Println(domain)
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
			result[0] = strings.ReplaceAll(result[0], "https://", "")
			result[0] = strings.Split(result[0], ":")[0]
			row := []string{result[0], result[1], result[2]}
			data = append(data, row)
		}

		table.SetFooter([]string{"", "Total", strconv.Itoa(len(domainList.Results))}) // Add Footer
		table.AppendBulk(data)                                                        // Add Bulk Data
		table.Render()
	}
}

func echoerr() {
	fmt.Println("[-] 调用FOFA-API失败！")
}
