//go:build windows

//link:https://github.com/eeeeeeeeee-code/e0e1-config/tree/main/pkg/navicat

package navicat

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/blowfish"
	"golang.org/x/sys/windows/registry"
	"strings"
)

var (
	aesKey        = []byte("libcckeylibcckey")
	aesIV         = []byte("libcciv libcciv ")
	blowfishKey   = sha1Sum([]byte("3DC5CA39"))
	blowfishIV    = hexDecode("d9c7c3c8870d64bd")
	resultBuilder strings.Builder
)

func Run(cmd *cobra.Command, args []string) {
	navinfo, err := ScanNavicat()
	if err != nil {
		fmt.Printf("获取Navicat信息失败! %v\n", err)
		return
	}
	fmt.Println(navinfo)
}

func sha1Sum(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}

func hexDecode(s string) []byte {
	decoded, _ := hex.DecodeString(s)
	return decoded
}

func decryptNavicat11(hexPassword string) (string, error) {
	if hexPassword == "" {
		return "无密码", nil
	}

	encryptedData, err := hex.DecodeString(strings.ToLower(hexPassword))
	if err != nil {
		return "", fmt.Errorf("十六进制解码失败: %v", err)
	}

	cipher, err := blowfish.NewCipher(blowfishKey)
	if err != nil {
		return "", fmt.Errorf("创建Blowfish密码器失败: %v", err)
	}

	roundCount := len(encryptedData) / 8
	decryptedPassword := make([]byte, 0)
	currentVector := make([]byte, len(blowfishIV))
	copy(currentVector, blowfishIV)

	for i := 0; i < roundCount; i++ {
		block := encryptedData[i*8 : (i+1)*8]
		decryptedBlock := make([]byte, 8)
		cipher.Decrypt(decryptedBlock, block)

		for j := 0; j < 8; j++ {
			decryptedBlock[j] ^= currentVector[j]
		}
		decryptedPassword = append(decryptedPassword, decryptedBlock...)

		for j := 0; j < 8; j++ {
			currentVector[j] ^= block[j]
		}
	}

	if len(encryptedData)%8 != 0 {
		lastBlock := make([]byte, 8)
		cipher.Encrypt(lastBlock, currentVector)

		for i := 0; i < len(encryptedData)%8; i++ {
			decryptedPassword = append(decryptedPassword, encryptedData[roundCount*8+i]^lastBlock[i])
		}
	}

	return strings.TrimRight(string(decryptedPassword), "\x00"), nil
}

func decryptNavicat12(hexPassword string) (string, error) {
	if hexPassword == "" {
		return "无密码", nil
	}

	encryptedData, err := hex.DecodeString(strings.ToLower(hexPassword))
	if err != nil {
		return "", fmt.Errorf("十六进制解码失败: %v", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", fmt.Errorf("创建AES密码器失败: %v", err)
	}

	mode := cipher.NewCBCDecrypter(block, aesIV)
	decryptedData := make([]byte, len(encryptedData))
	mode.CryptBlocks(decryptedData, encryptedData)

	return strings.TrimRight(string(decryptedData), "\x00"), nil
}

func DecryptPassword(encryptedPassword string, version int) string {
	if encryptedPassword == "" {
		return "无密码"
	}

	var result string
	var err error

	if version == 11 {
		result, err = decryptNavicat11(encryptedPassword)
	} else if version >= 12 {
		result, err = decryptNavicat12(encryptedPassword)
	} else {
		return "[-] 不支持的版本"
	}

	if err != nil {
		return fmt.Sprintf("[-] 解密失败: %v", err)
	}

	return strings.TrimSpace(result)
}

func GetNavicatServers() ([]string, error) {
	baseKey := `Software\PremiumSoft`
	var allConnections []string

	key, err := registry.OpenKey(registry.CURRENT_USER, baseKey, registry.READ)
	if err != nil {
		return nil, fmt.Errorf("打开注册表项失败: %v", err)
	}
	defer key.Close()

	subKeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return nil, fmt.Errorf("读取子键失败: %v", err)
	}

	for _, subKey := range subKeys {
		if !strings.Contains(subKey, "Navicat") {
			continue
		}

		serverPath := fmt.Sprintf("%s\\%s\\Servers", baseKey, subKey)
		serverKey, err := registry.OpenKey(registry.CURRENT_USER, serverPath, registry.READ)
		if err != nil {

			continue
		}

		serverNames, err := serverKey.ReadSubKeyNames(-1)

		if err != nil {
			serverKey.Close()
			continue
		}

		for _, serverName := range serverNames {
			connection, err := getServerInfo(serverPath, serverName)
			if err == nil {
				allConnections = append(allConnections, connection...)
			}
		}

		serverKey.Close()
	}

	return allConnections, nil
}

func getServerInfo(serverPath, serverName string) ([]string, error) {
	fullPath := fmt.Sprintf("%s\\%s", serverPath, serverName)
	key, err := registry.OpenKey(registry.CURRENT_USER, fullPath, registry.READ)
	if err != nil {
		return []string{}, err
	}
	defer key.Close()

	valueNames, err := key.ReadValueNames(-1)
	values := []string{}
	if err != nil {
		fmt.Printf("读取值名称失败: %v\n", err)
	} else {
		values = append(values, fmt.Sprintf("连接 %s 包含以下字段:\n", serverName))
		for _, name := range valueNames {
			value, _, err := key.GetStringValue(name)
			if err == nil && value != "" {

				if name == "Password" || name == "Pwd" {
					decryptedValue := DecryptPassword(value, 11)
					values = append(values, fmt.Sprintf("  %s: %s (解密后: %s)\n", name, value, decryptedValue))
				} else {
					values = append(values, fmt.Sprintf("  %s: %s\n", name, value))
				}
			} else if err == nil && value == "" {

				continue
			} else {

				binValue, _, err := key.GetBinaryValue(name)
				if err == nil && len(binValue) > 0 {
					fmt.Printf("  %s: [二进制数据，长度: %d]\n", name, len(binValue))
				} else {

					intValue, _, err := key.GetIntegerValue(name)
					if err == nil {
						values = append(values, fmt.Sprintf("  %s: %d\n", name, intValue))
					}
				}
			}
		}
	}

	return values, nil
}

func ScanNavicat() (string, error) {
	connections, err := GetNavicatServers()
	if err != nil {
		return "", fmt.Errorf("从注册表获取Navicat连接失败: %v", err)
	}

	if len(connections) > 0 {
		resultBuilder.WriteString("[+] 成功从注册表获取保存的 Navicat 连接\n")

		resultBuilder.WriteString(fmt.Sprintf("[+] %s\n", strings.Join(connections, "")))

	} else {
		return "", fmt.Errorf("未找到任何包含密码的 Navicat 连接")
	}
	return resultBuilder.String(), nil
}
