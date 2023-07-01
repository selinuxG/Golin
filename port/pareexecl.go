package port

import "strings"

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
