package web

import (
	"bufio"
	"encoding/json"
	"fmt"
	"golin/global"
	"os"
)

type Service struct {
	Name   string `json:"name"`
	Ip     string `json:"ip"`
	User   string `json:"user"`
	Port   string `json:"port"`
	Type   string `json:"type"`
	Time   string `json:"time"`
	Status string `json:"status"`
}

const (
	Success = "成功"
	Failed  = "失败"
)

// WriteJSONToHistory 将结构体切片转换为JSON并写入文件
func WriteJSONToHistory(service Service) {
	data, _ := json.Marshal(service)
	writedata := string(data) + "\n"
	global.AppendToFile(global.Succwebpath, writedata)
}

func parseJSONFile() ([]Service, error) {
	file, err := os.Open(global.Succwebpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var services []Service

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		var service Service
		err := json.Unmarshal(line, &service)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

// ReadJSONFromFile 文件中读取JSON并将其转换为结构体切片
func ReadJSONFromFile() ([]Service, error) {
	data, err := os.ReadFile(global.Succwebpath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var allserver []Service
	err = json.Unmarshal(data, &allserver)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return allserver, nil
}
