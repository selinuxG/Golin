//go:build windows

package windows

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
)

// regedt 读取是user 计算机\HKEY_LOCAL_MACHINE
func regedt(path, keys string) (string, error) {
	// 打开注册表键值
	key, err := registry.OpenKey(registry.CURRENT_USER, path, registry.READ)
	if err != nil {
		fmt.Println("Error opening registry key:", err)
		return "", err
	}
	defer key.Close()
	// 读取特定键
	value, _, err := key.GetStringValue(keys)
	if err != nil {
		return "", err
	}

	return value, nil
}

// regedtmachine 读取的是计算机\HKEY_LOCAL_MACHINE\
func regedtmachineDWORD(path, keys string) (uint64, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE)
	if err != nil {
		fmt.Println("Error opening registry key:", err)
		return 0, err
	}
	defer key.Close()
	// 读取特定键

	value, _, err := key.GetIntegerValue(keys)
	if err != nil {
		return 0, err
	}

	return value, nil
}

// regedtmachineString 读取的是计算机\HKEY_LOCAL_MACHINE\
func regedtmachineString(path, keys string) (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE)
	if err != nil {
		fmt.Println("Error opening registry key:", err)
		return "", err
	}
	defer key.Close()
	// 读取特定键

	value, _, err := key.GetStringValue(keys)
	if err != nil {
		return "", err
	}

	return value, nil
}
