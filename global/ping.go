package global

import (
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// NetWorkStatus 检查ping 返回是否可ping通以及操作系统
func NetWorkStatus(ip string) (bool, string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "4", ip)
	} else {
		cmd = exec.Command("ping", "-c", "4", ip)
	}

	output, err := cmd.Output()
	if err != nil {
		return false, ""
	}
	outttl := strings.ToLower(string(output)) //所有大写转换为小写
	if strings.Contains(outttl, "ttl") {
		// Extract TTL value
		re := regexp.MustCompile(`ttl=(\d+)`)
		ttlStr := re.FindStringSubmatch(outttl)
		if len(ttlStr) > 1 {
			ttl, _ := strconv.Atoi(ttlStr[1])
			switch {
			case ttl <= 64:
				return true, "Linux/Unix"
			case ttl <= 128:
				return true, "Windows"
			default:
				return true, "Unknown"
			}
		}
	}
	return false, ""
}
