package run

import (
	"fmt"
	"go.uber.org/zap"
	"golin/config"
	"golin/global"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var (
	succpath     = global.Succpath     //保存采集目录
	Split        = global.Split        //默认分隔符
	count        int                   //总数量,多少行文件就是多少
	wg           sync.WaitGroup        //线程
	errhost      []string              //失败主机列表
	runcmd       = ""                  //运行的linux命令
	denynametype = global.Denynametype //windos下不允许创建名称的特殊符号。
)

var (
	zlog = config.Log //自定义日志
)

// Rangefile 遍历文件并创建线程 path=模式目录 spr=按照什么分割 runtype运行类型
func Rangefile(path string, spr string, runtype string) {
	fire, _ := os.ReadFile(path)
	lines := strings.Split(string(fire), "\n")
	wg.Add(len(lines))
	for i := 0; i < len(lines); i++ {
		//如果是空行则跳过线程减1
		if len(lines[i]) == 0 {
			wg.Done()
			continue
		}
		firecount := strings.Count(lines[i], spr)
		if firecount != 4 {
			wg.Done()
			zlog.Warn("主机格式不正确，跳过！")
			continue
		}
		//总数量+1
		count += 1
		linedata := lines[i]
		Name := strings.Split(string(linedata), spr)[0]
		Host := strings.Split(string(linedata), spr)[1]
		User := strings.Split(string(linedata), spr)[2]
		Passwrod := strings.Split(string(linedata), spr)[3]
		Port1 := strings.Split(string(linedata), spr)[4]
		//windos中换行符可能存在为/r/n,之前分割/n,还留存/r,清除它
		Porttmp := strings.Replace(Port1, "\r", "", -1)
		Port, err := strconv.Atoi(Porttmp)
		if err != nil {
			wg.Done()
			zlog.Warn("端口转换失败: ", zap.String("IP", Host))
			errhost = append(errhost, Host)
			continue
		}
		//判断host是不是正确的IP地址格式
		address := net.ParseIP(Host)
		if address == nil {
			wg.Done()
			zlog.Warn("IP地址格式不正确，跳过！", zap.String("IP", Host))
			count = count - 1
			continue
		}
		//判断端口范围是否是1-65535
		if Port == 0 || Port > 65535 {
			wg.Done()
			zlog.Warn("端口范围不正确，跳过！", zap.String("IP", Host), zap.Int("Port:", Port))
			count = count - 1
			continue
		}
		//如果是Windows先判断保存文件是否存在特殊字符,是的话不执行直接记录为失败主机
		if runtime.GOOS == "windows" {
			if InSlice(denynametype, Name) {
				wg.Done()
				zlog.Warn("名称存在特殊符号，跳过！")
				errhost = append(errhost, Host)
				continue
			}
		}

		firepath := filepath.Join(succpath, runtype)
		_, err = os.Stat(firepath)
		if err != nil {
			_ = os.MkdirAll(firepath, os.FileMode(global.FilePer))
		}

		switch runtype {
		case "Linux":
			go func() {
				sshErr := Runssh(Name, Host, User, Passwrod, Port, runcmd)
				if sshErr != nil {
					errhost = append(errhost, Host)
					zlog.Warn("采集Linux安全配置失败:", zap.Error(sshErr))
				}
			}()
		case "Mysql":
			go RunMysql(Name, User, Passwrod, Host, strconv.Itoa(Port))
		case "Redis":
			go Runredis(Name, User, Host, Passwrod, strconv.Itoa(Port))
		case "pgsql":
			go Pgsql(Name, Host, User, Passwrod, strconv.Itoa(Port))
		case "sqlserver":
			go SqlServerrun(Name, Host, User, Passwrod, strconv.Itoa(Port))
		case "oracle":
			go func() {
				Err := OracleRun(Name, Host, User, Passwrod, strconv.Itoa(Port))
				if Err != nil {
					errhost = append(errhost, Host)
					zlog.Warn("采集Oracle安全配置失败:", zap.Error(Err))
				}
			}()
		}
	}
	wg.Wait()
}

// Onlyonerun 只允许一次的模式
func Onlyonerun(value string, spr string, runtype string) {
	firecount := strings.Count(value, spr)
	if firecount != 4 {
		zlog.Warn("错误！格式不正确，退出！（默认为：名称~IP~用户~密码~端口）")
		return
	}
	Name := strings.Split(value, spr)[0]
	Host := strings.Split(value, spr)[1]
	User := strings.Split(value, spr)[2]
	Passwrod := strings.Split(value, spr)[3]
	Port1 := strings.Split(value, spr)[4]
	Porttmp := strings.Replace(Port1, "\r", "", -1)
	Port, err := strconv.Atoi(Porttmp)
	if err != nil {
		zlog.Warn("错误！端口格式转换失败,退出！")
		return
	}
	address := net.ParseIP(Host)
	if address == nil {
		zlog.Warn("错误！不是正确的IP地址,退出！")
		return
	}
	//判断端口范围是否是1-65535
	if Port == 0 || Port > 65535 {
		zlog.Warn("错误！不是正确的端口范围,退出！")
		return
	}
	//如果是Windows先判断保存文件是否存在特殊字符,是的话不执行直接记录为失败主机
	if runtime.GOOS == "windows" {
		if InSlice(denynametype, Name) {
			zlog.Warn("错误！保存文件包含特殊字符,无法保存,请修改在执行！")
			return
		}
	}

	firepath := filepath.Join(succpath, runtype)
	_, err = os.Stat(firepath)
	if err != nil {
		_ = os.MkdirAll(firepath, os.FileMode(global.FilePer))
	}

	switch runtype {
	case "Linux":
		wg.Add(1)
		go func() {
			sshErr := Runssh(Name, Host, User, Passwrod, Port, runcmd)
			if sshErr != nil {
				errhost = append(errhost, Host)
				zlog.Warn("采集Linux安全配置失败:", zap.Error(sshErr))
			}
		}()
	case "MySQL":
		wg.Add(1)
		config.Log.Info("开启运行Mysql模式", zap.String("名称:", Name), zap.String("IP", Host))
		go RunMysql(Name, User, Passwrod, Host, strconv.Itoa(Port))
	case "Redis":
		wg.Add(1)
		config.Log.Info("开启运行Redis模式", zap.String("名称:", Name), zap.String("IP", Host))
		go Runredis(Name, User, Host, Passwrod, strconv.Itoa(Port))
	case "pgsql":
		wg.Add(1)
		config.Log.Info("开启运行pgsql模式", zap.String("名称:", Name), zap.String("IP", Host))
		go Pgsql(Name, Host, User, Passwrod, strconv.Itoa(Port))
	case "sqlserver":
		wg.Add(1)
		config.Log.Info("开启运行sqlserver模式", zap.String("名称:", Name), zap.String("IP", Host))
		go SqlServerrun(Name, Host, User, Passwrod, strconv.Itoa(Port))
	case "oracle":
		wg.Add(1)
		config.Log.Info("开启运行oracle模式", zap.String("名称:", Name), zap.String("IP", Host))
		go func() {
			Err := OracleRun(Name, Host, User, Passwrod, strconv.Itoa(Port))
			if Err != nil {
				errhost = append(errhost, Host)
				zlog.Warn("采集Oracle安全配置失败:", zap.Error(Err))
			}
		}()
	}
	wg.Wait() //等待运行结束
}

// Checkfile 判断某个模式下的默认文件是否存在
func Checkfile(name string, data string, pems int, path string) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		datanew := []byte(string(data))
		os.WriteFile(path, datanew, fs.FileMode(pems))
		config.Log.Warn("默认文件不存在", zap.String("默认文件", name))
		config.Log.Info("已自动创建符合格式的默认文件，修改后再来吧！", zap.String("默认文件", name))
		os.Exit(3)
	}
}

// Deffile 程序退出前运行的函数，用于生成日志
func Deffile(moude string, count int, success int, errhost []string) {
	//如果是Windows操作系统则自动弹出采集完成文件夹
	//defer func() {
	//	switch localos := runtime.GOOS; localos {
	//	case "windows":
	//		_, err := os.Stat(succpath)
	//		if err == nil {
	//			if success > 0 {
	//				cmd := exec.Command("powershell", "Get-Process | Where-Object {$_.MainWindowTitle -like '*"+succpath+"*'}") // 检查该目录是否已经打开
	//				output, _ := cmd.Output()
	//				if len(output) == 0 { // 未打开目录
	//					cmd = exec.Command("explorer.exe", succpath)
	//					err := cmd.Start()
	//					if err != nil {
	//						panic(err)
	//					}
	//				}
	//			}
	//		}
	//	}
	//}()
	if count == success {
		zlog.Info("运行记录",
			zap.String("执行模式", moude),
			zap.Int("采集总数量", count),
			zap.Int("采集成功数量", success),
			zap.String("成功率", "100%"),
		)
		return
	}
	zlog.Warn("运行记录",
		zap.String("执行模式", moude),
		zap.Int("采集总数量:", count),
		zap.Int("采集成功数量:", success),
		zap.Int("采集失败数量:", count-success),
		zap.String("成功率", fmt.Sprintf("%.2f%%", calculateSuccessRate(success, count))),
	)
	if len(errhost) > 0 {
		for _, v := range errhost {
			config.Log.Warn("失败记录", zap.String("运行模式", moude), zap.String("IP", v))
		}
	}
}

// InSlice 判断字符串是否在 不允许命名的slice中。
func InSlice(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// calculateSuccessRate 计算成功率
func calculateSuccessRate(success, count int) float64 {
	successFloat := float64(success)
	countFloat := float64(count)
	successRate := (successFloat / countFloat) * 100.0
	return successRate
}
