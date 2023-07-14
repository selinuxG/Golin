package port

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func saveXlsx(infoList []INFO, ipList []string) {

	if len(infolist) == 0 && len(ipList) == 0 {
		return
	}

	f := excelize.NewFile()
	defer f.Close()
	// 创建一个工作表
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		fmt.Printf("[-] 保存文件失败！1、文件是否已打开？ 2、是否有权限\n")
		return
	}

	//设置列宽度，A为序号调小一点
	_ = f.SetColWidth("Sheet1", "A", "B", 8)
	_ = f.SetColWidth("Sheet1", "B", "C", 20)
	_ = f.SetColWidth("Sheet1", "C", "D", 10)
	_ = f.SetColWidth("Sheet1", "D", "E", 80)

	// 设置单元格的值
	_ = f.SetCellValue("Sheet1", "A1", "序号")
	_ = f.SetCellValue("Sheet1", "B1", "主机")
	_ = f.SetCellValue("Sheet1", "C1", "端口")
	_ = f.SetCellValue("Sheet1", "D1", "协议及组件")
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	cell := 1

	for _, info := range infoList {
		cell += 1
		_ = f.SetCellValue("Sheet1", fmt.Sprintf("A%d", cell), cell-1)
		_ = f.SetCellValue("Sheet1", fmt.Sprintf("B%d", cell), info.Host)
		_ = f.SetCellValue("Sheet1", fmt.Sprintf("C%d", cell), info.Port)
		_ = f.SetCellValue("Sheet1", fmt.Sprintf("D%d", cell), info.Protocol)
	}

	//存活的主机
	index, err = f.NewSheet("Sheet2")
	if err != nil {
		fmt.Printf("[-] 保存文件失败！1、文件是否已打开？ 2、是否有权限\n")
		return
	}
	_ = f.SetColWidth("Sheet2", "A", "B", 8)
	_ = f.SetColWidth("Sheet2", "B", "C", 20)
	_ = f.SetCellValue("Sheet2", "A1", "序号")
	_ = f.SetCellValue("Sheet2", "B1", "主机")
	cell = 1
	for _, ip := range ipList {
		cell += 1
		_ = f.SetCellValue("Sheet2", fmt.Sprintf("A%d", cell), cell-1)
		_ = f.SetCellValue("Sheet2", fmt.Sprintf("B%d", cell), ip)
	}

	// 根据指定路径保存文件
	if err := f.SaveAs("portscan.xlsx"); err != nil {
		fmt.Printf("[-] 保存文件失败！1、文件是否已打开？ 2、是否有权限\n")
		return
	}

	fmt.Printf("[*] 结果保存路径：\033[32m%s\033[0m\n", "portscan.xlsx")

}
