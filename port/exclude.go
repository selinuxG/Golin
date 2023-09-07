package port

import (
	"golin/global"
	"os"
	"strings"
)

// excludePort 排除端口
func excludePort(port string) {
	if port == "" {
		portlist = removeDuplicates(portlist) //现有端口去重
		return
	}

	if strings.Count(port, ",") == 0 {
		port = port + ","
	}

	for _, v := range strings.Split(port, ",") {
		for i := 0; i < len(portlist); i++ {
			if portlist[i] == v {
				portlist = append(portlist[:i], portlist[i+1:]...) // 删除与 v 相等的元素
				i--                                                // 更新索引，因为切片长度减少了 1
			}
		}
	}

}

// removeDuplicates 切片去重
func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var list []string

	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// removeIP 获取不进行扫描的IP
func removeIP(data string) {
	if data == "" {
		return
	}

	var NoList []string

	if global.PathExists(data) {
		fireData, _ := os.ReadFile(data)
		fireData = []byte(strings.ReplaceAll(string(fireData), "\r\n", "\n"))
		for _, ip := range strings.Split(string(fireData), "\n") {
			NoList = append(NoList, ip)
		}
	}

	// 创建一个映射，用于快速查找nolist中的元素
	noMap := make(map[string]bool)
	for _, ip := range NoList {
		noMap[ip] = true
	}
	// 创建一个新的切片，其中只包含那些不在nolist中的iplist元素
	var NewList []string
	for _, ip := range iplist {
		if _, ok := noMap[ip]; !ok {
			NewList = append(NewList, ip)
		}
	}
	iplist = NewList

}
