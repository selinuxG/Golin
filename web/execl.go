package web

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"golin/global"
)

// CreateTmpTxt 生成xlsx转txt文件
func CreateTmpTxt(xlsx, txt string) bool {
	f, err := excelize.OpenFile(xlsx)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer f.Close()
	// 获取 Sheet1 上所有单元格
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return false
	}
	for v, row := range rows {
		iplist := []string{}
		for _, colCell := range row {
			if v == 0 {
				continue
			}
			iplist = append(iplist, colCell)
		}
		if len(iplist) > 0 {
			err := global.AppendToFile(txt, fmt.Sprintf("%v~~%v~~%v~~%v~~%v\n", iplist[0], iplist[1], iplist[2], iplist[3], iplist[4]))
			if err != nil {
				return false
			}
		}
	}

	return true

}

// CreateTemplateXlsx 生成golin上传模板文件
func CreateTemplateXlsx() bool {
	f := excelize.NewFile()
	defer f.Close()
	// 创建一个工作表
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		fmt.Println(err)
		return false
	}
	// 设置单元格的值
	f.SetCellValue("Sheet1", "A1", "名称")
	f.SetCellValue("Sheet1", "B1", "IP")
	f.SetCellValue("Sheet1", "C1", "用户")
	f.SetCellValue("Sheet1", "D1", "密码")
	f.SetCellValue("Sheet1", "E1", "端口")
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	// 根据指定路径保存文件
	if err := f.SaveAs(global.XlsxTemplateName); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
