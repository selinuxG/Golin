package scan

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"golin/poc"
	"golin/scan/crack"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

var (
	dir = "ScanLog" //扫描结果保存目录
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

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println("[-] 创建ScanLog目录失败!")
			return
		}
	}

	filename := time.Now().Format("200601021504") + "-GoLin.xlsx"
	filename = filepath.Join("ScanLog", filename)
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
	if len(crack.MapCrackHost) <= 0 {
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

	cell := 1
	for _, sussCrack := range crack.MapCrackHost {
		cell += 1
		setCellValues(f, sheet, cell, []interface{}{cell - 1, sussCrack.Host, strconv.Itoa(sussCrack.Port), sussCrack.User, sussCrack.Passwd, sussCrack.Mode})
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
