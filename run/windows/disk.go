//go:build windows

package windows

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func disk() {
	value := ExecCommands(`wmic logicaldisk get DeviceID, FileSystem, Size, FreeSpace`)
	value = strings.ReplaceAll(value, "\r\r\n", "\n")
	echo := ""
	for _, v := range strings.Split(value, "\n") {
		if strings.Count(v, "DeviceID") > 0 {
			continue
		}
		re := regexp.MustCompile(`\s{2,}`)
		v = re.ReplaceAllString(v, " ")
		data := strings.Split(v, " ")
		if len(data) == 5 {
			free, _ := strconv.ParseUint(data[2], 10, 64) //如果为int会溢出,int类型在32位系统上的最大值是2147483647
			freeg := fmt.Sprintf("%d/G", free/1024.0/1024.0/1024.0)
			all, _ := strconv.ParseUint(data[3], 10, 64) //如果为int会溢出,int类型在32位系统上的最大值是2147483647
			allg := fmt.Sprintf("%d/G", all/1024.0/1024.0/1024.0)
			if freeg == "0/G" || allg == "0/G" {
				freeg = fmt.Sprintf("%d/M", free/1024.0/1024.0)
				allg = fmt.Sprintf("%d/M", all/1024.0/1024.0)
			}
			percentage := (float64(free) / float64(all)) * 100
			percentageStr := fmt.Sprintf("%.2f%%", percentage)

			echo += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", data[0], data[1], freeg, allg, percentageStr)
		}
	}
	html = strings.ReplaceAll(html, "磁盘信息结果", echo)
}
