package global

import (
	"os/exec"
	"runtime"
	"strings"
)

// NetWorkStatus 检查ping
func NetWorkStatus(ip string) bool {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "2", ip)
	} else {
		cmd = exec.Command("ping", "-c", "2", ip)
	}

	output, err := cmd.Output()
	if err != nil {
		return false
	}
	outttl := strings.ToLower(string(output)) //所有大写转换为小写
	if strings.Contains(outttl, "ttl") {
		return true
	}
	return false
}
