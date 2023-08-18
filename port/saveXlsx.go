package port

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"golin/poc"
	"golin/port/crack"
	"regexp"
	"time"
)

func saveXlsx(infoList []INFO, ipList []string) {
	if len(infoList) == 0 && len(ipList) == 0 {
		return
	}

	f := excelize.NewFile()
	f.SetSheetName("Sheet1", "扫描端口结果")
	defer f.Close()

	_ = createInfoSheet(f, "扫描端口结果", infoList)
	_ = createIpSheet(f, "存活主机", ipList)
	_ = createPoc(f, "漏洞资产")
	_ = createCrack(f, "弱口令资产")

	filename := time.Now().Format("20060102") + "-GoLin.xlsx"
	if err := f.SaveAs(filename); err != nil {
		fmt.Printf("[-] 保存文件失败！1、文件是否已打开？ 2、是否有权限\n")
		fmt.Println(err)
		return
	}

	fmt.Printf(clearLine)
	fmt.Printf("[*] 结果保存路径：\033[32m%s\033[0m\n", filename)
}

// createInfoSheet 创建端口开放信息的sheet
func createInfoSheet(f *excelize.File, sheet string, infoList []INFO) error {
	if len(infoList) == 0 {
		return errors.New("开放端口为空跳过")
	}
	index, err := f.NewSheet(sheet)
	if err != nil {
		return err
	}

	setColWidths(f, sheet)
	f.SetColWidth(sheet, "D", "E", 100)
	setHeaderValues(f, sheet, []string{"序号", "主机", "端口", "协议及组件"})

	f.SetActiveSheet(index)

	for i, info := range infoList {
		cell := i + 2
		setCellValues(f, sheet, cell, []interface{}{i + 1, info.Host, info.Port, cleanProtocol(info.Protocol)})
	}

	return nil
}

// createIpSheet 创建存活主机的sheet
func createIpSheet(f *excelize.File, sheet string, ipList []string) error {
	if len(infolist) == 0 {
		return errors.New("存活主机为空跳过")
	}
	_, err := f.NewSheet(sheet)
	if err != nil {
		return err
	}

	setColWidths(f, sheet)
	setHeaderValues(f, sheet, []string{"序号", "主机"})

	for i, ip := range ipList {
		cell := i + 2
		setCellValues(f, sheet, cell, []interface{}{i + 1, ip})
	}

	return nil
}

// createCrack 创建弱口令sheet
func createCrack(f *excelize.File, sheet string) error {
	if len(crack.ListCrackHost) <= 0 {
		return errors.New("弱口令资产为空")
	}
	_, err := f.NewSheet(sheet)
	if err != nil {
		return err
	}

	setColWidths(f, sheet)
	f.SetColWidth(sheet, "D", "E", 10)
	f.SetColWidth(sheet, "E", "F", 10)
	f.SetColWidth(sheet, "F", "G", 10)

	setHeaderValues(f, sheet, []string{"序号", "主机", "端口", "用户", "密码", "模式"})

	for i, v := range crack.ListCrackHost {
		cell := i + 2
		setCellValues(f, sheet, cell, []interface{}{i + 1, v.Host, v.Port, v.User, v.Passwd, v.Mode})
	}

	return nil
}

// createPoc 创建漏洞资产sheet
func createPoc(f *excelize.File, sheet string) error {
	if len(poc.ListPocInfo) == 0 {
		return errors.New("漏洞资产为空")
	}
	_, err := f.NewSheet(sheet)
	if err != nil {
		return err
	}
	f.SetColWidth(sheet, "A", "B", 8)
	f.SetColWidth(sheet, "B", "C", 40)
	f.SetColWidth(sheet, "C", "D", 40)
	f.SetColWidth(sheet, "D", "E", 40)

	setHeaderValues(f, sheet, []string{"序号", "URL", "漏洞POC", "信息描述"})

	for i, v := range poc.ListPocInfo {
		cell := i + 2
		setCellValues(f, sheet, cell, []interface{}{i + 1, v.Url, v.Cve, v.Flag})
	}

	return nil
}

func setColWidths(f *excelize.File, sheet string) {
	f.SetColWidth(sheet, "A", "B", 8)
	f.SetColWidth(sheet, "B", "C", 20)
	f.SetColWidth(sheet, "C", "D", 10)
}

func setHeaderValues(f *excelize.File, sheet string, headers []string) {
	for i, header := range headers {
		colLetter := string('A' + i)
		f.SetCellValue(sheet, colLetter+"1", header)
	}
}

func setCellValues(f *excelize.File, sheet string, cell int, values []interface{}) {
	for i, value := range values {
		colLetter := string('A' + i)
		f.SetCellValue(sheet, fmt.Sprintf("%s%d", colLetter, cell), value)
	}
}

func cleanProtocol(protocol string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(protocol, "")
}
