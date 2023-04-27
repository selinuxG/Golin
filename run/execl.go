package run

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
	"golin/global"
	"os"
	"strconv"
	"strings"
)

type Iplist struct {
	Name   string
	IP     string
	Port   string
	User   string
	Passwd string
}

var (
	listIp   []Iplist
	modeip   = make(map[string]string)
	modetype = []string{"name", "ip", "port", "user", "passwd"}
)

func Execl(cmd *cobra.Command, args []string) {
	filename, _ := cmd.Flags().GetString("file")
	ip, _ := cmd.Flags().GetString("ip")
	name, _ := cmd.Flags().GetString("name")
	port, _ := cmd.Flags().GetString("port")
	user, _ := cmd.Flags().GetString("user")
	passwd, _ := cmd.Flags().GetString("passwd")
	sheet, _ := cmd.Flags().GetString("sheet")
	savafile, _ := cmd.Flags().GetString("sava")
	if filename == "" && ip == "" && user == "" && passwd == "" {
		fmt.Println("err: 字段值不可为空")
		return
	}
	modeip["name"] = name
	modeip["ip"] = ip
	modeip["port"] = port
	modeip["user"] = user
	modeip["passwd"] = passwd
	readColumnFromExcel(filename, sheet, modeip)
	os.Remove(savafile)
	file, err := os.OpenFile(savafile, os.O_CREATE|os.O_APPEND, os.FileMode(global.FilePer))
	if err != nil {
		fmt.Println("打开文件失败,", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(file)
	write := bufio.NewWriter(file)
	for _, v := range listIp {
		_, err := write.WriteString(fmt.Sprintf("%s~%s~%s~%s~%s\n", v.Name, v.IP, v.User, v.Passwd, v.Port))
		if err != nil {
			fmt.Println("写入文件失败...", err)
			return
		}
	}
	//文件写入
	fmt.Printf("结果保存在%s,正在保存中...\n", savafile)
	err = write.Flush()
	if err != nil {
		fmt.Println(err)
		return
	}
}

// readColumnFromExcel 基于xlsx生成txt文件
func readColumnFromExcel(filename, sheetName string, columnName map[string]string) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println(filename, err)
		os.Exit(3)
	}
	// 获取工作表的行数
	rows, err := f.GetRows(sheetName)
	if err != nil {
		fmt.Println(err)
		return
	}
	newip := Iplist{}
	for i := 1; i <= len(rows); i++ {
		for _, v := range modetype {
			cellCoord := fmt.Sprintf("%s%d", columnName[v], i)       // 设置单元格的坐标，例如 "B1", "B2", "B3" 等
			cellValuename, _ := f.GetCellValue(sheetName, cellCoord) //根据坐标获取该单元格的值
			switch v {
			case "ip":
				if strings.Count(cellValuename, ":") == 1 {
					newip.IP = strings.Split(cellValuename, ":")[0]   //ip
					newip.Port = strings.Split(cellValuename, ":")[1] //port
					continue
				}
				newip.IP = cellValuename
			case "name":
				newip.Name = cellValuename
			case "user":
				newip.User = cellValuename
			case "passwd":
				newip.Passwd = cellValuename
			case "port":
				if newip.Port == "" {
					port, err := strconv.Atoi(cellValuename)
					if err != nil {
						continue
					}
					if port > 0 && port < 65535 {
						newip.Port = cellValuename
					}
				}
			}
		}
		if newip.IP != "" && newip.User != "" && newip.Passwd != "" && newip.Port != "" {
			listIp = append(listIp, newip)
		}
	}
}
