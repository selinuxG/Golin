package scan

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"golin/global"
	"golin/poc"
	"golin/scan/crack"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
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

	filename := filepath.Join("ScanLog", global.XlsxFileName)
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
	setHeaderValues(f, sheet, []string{"序号", "主机", "操作系统", "资产标签"})

	for i, ip := range ipList {
		cell := i + 2
		setCellValues(f, sheet, cell, []interface{}{
			i + 1,
			ip,
			global.LoadOrDefault(&IPListOS, ip, "未知"),
			TagAsset(ip),
		})
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
		colLetter, _ := excelize.ColumnNumberToName(i + 1) // Excel列是从 1 开始的
		f.SetCellValue(sheet, colLetter+"1", header)
	}
}

func setCellValues(f *excelize.File, sheet string, cell int, values []interface{}) {
	for i, value := range values {
		colLetter, _ := excelize.ColumnNumberToName(i + 1)
		cellRef := fmt.Sprintf("%s%d", colLetter, cell)
		f.SetCellValue(sheet, cellRef, value)
	}
}

func cleanProtocol(protocol string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(protocol, "")
}

// TagAsset 资产打标签
func TagAsset(ip string) string {
	// 将服务类型映射到端口号切片
	servicePorts := map[string][]int{
		"WEB服务器":  {80, 443, 8080},
		"数据库服务器":  {3306, 1433, 5432, 1521, 5236},
		"缓存服务器":   {6379, 27017},
		"FTP服务器":  {21},
		"VPN服务器":  {1194, 1723},
		"监控服务器":   {9090, 3000}, // Prometheus (9090), Grafana (3000)
		"消息队列服务器": {5672, 9092}, // RabbitMQ (5672), Kafka (9092)
		"搜索引擎服务器": {8983, 9200}, // Solr (8983), Elasticsearch (9200)
		"版本控制服务器": {9418, 3690}, // Git (9418), Subversion (3690)
	}

	// 使用map来去重和存储标签
	var tagList = make(map[string]bool)

	listenPorts := GetPortsByHost(infolist, ip)
	// 遍历监听的端口
	for _, portInfo := range listenPorts {
		port, err := strconv.Atoi(portInfo)
		if err != nil {
			continue // 如果端口号不是有效的整数，跳过
		}

		// 检查端口属于哪个服务类型
		for service, ports := range servicePorts {
			for _, servicePort := range ports {
				if port == servicePort {
					tagList[service] = true
					break
				}
			}
		}
	}

	// 如果没有任何标签，也标记为应用服务器
	if len(tagList) == 0 {
		return ""
	}
	var tagListSlice []string
	for tag := range tagList {
		tagListSlice = append(tagListSlice, tag)
	}

	// 使用逗号分隔标签，创建一个单一字符串
	return strings.Join(tagListSlice, ",")
}

// GetPortsByHost 根据主机返回所有匹配的端口
func GetPortsByHost(infolist []INFO, host string) []string {
	var ports []string
	for _, info := range infolist {
		if info.Host == host {
			ports = append(ports, info.Port)
		}
	}
	return ports
}
