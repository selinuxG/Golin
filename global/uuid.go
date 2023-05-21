package global

import (
	"fmt"
	"github.com/google/uuid"
	"net"
	"strings"
)

// MacSearch 返回一个本机网络接口的MAC地址
func MacSearch() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("获取网络接口列表时出错:", err)
		return "", err
	}
	for _, inter := range interfaces {
		// 查找一个有效的MAC地址
		if inter.HardwareAddr != nil && len(strings.TrimSpace(inter.HardwareAddr.String())) > 0 {
			macAddress := inter.HardwareAddr.String()
			return macAddress, nil
		}
	}
	return "", err

}

// UuidFromMAC 根据MAC地址获取UUID
func UuidFromMAC(macAddress string) (string, error) {
	mac, err := net.ParseMAC(macAddress)
	if err != nil {
		return "", fmt.Errorf("解析MAC地址时出错: %v", err)
	}

	// 使用mac地址创建UUIDv5
	namespace := uuid.Nil
	deviceUUID := uuid.NewSHA1(namespace, mac)
	return deviceUUID.String(), nil
}
