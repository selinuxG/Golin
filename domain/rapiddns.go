package domain

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/olekukonko/tablewriter"
	"github.com/parnurzeal/gorequest"
	"os"
	"strconv"
	"strings"
	"time"
)

func rapidDNS(domain string) {

	request := gorequest.New().Timeout(10 * time.Second)
	if request == nil {
		return
	}
	domain = fmt.Sprintf("https://rapiddns.io/s/%s#result", domain)
	req := request.Get(domain)
	if req == nil {
		return
	}

	resp, body, errs := req.End()
	if len(errs) > 0 {
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return
	}
	var data [][]string
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"", "doamin", "ip", "type", "date"})

	doc.Find("tbody").Each(func(i int, s *goquery.Selection) {
		s.Find("tr").Each(func(j int, row *goquery.Selection) {
			var cells []string
			row.Find("td, th").Each(func(k int, cell *goquery.Selection) {
				cells = append(cells, strings.TrimSpace(cell.Text()))
			})
			rowdata := []string{cells[0], cells[1], cells[2], cells[3], cells[4]}
			data = append(data, rowdata)
		})
	})
	table.SetFooter([]string{"", "", "", "Total", strconv.Itoa(len(data))}) // Add Footer
	table.AppendBulk(data)                                                  // Add Bulk Data
	table.Render()

}
