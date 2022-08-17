package dbcp

import (
	_ "embed"
	"golin/dbcp/xml"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	//go:embed 07.安全计算环境-服务器和终端-高业尚服务器.xml
	gysxml    string
	files     []string
	writedata string
)

func Dirs() {
	f, err := ioutil.ReadDir("采集完成目录")
	if err != nil {
		log.Println(err)
		return
	}
	for _, file := range f {
		if findlinux := strings.Contains(file.Name(), "---linux.log"); findlinux {
			filepath := "采集完成目录\\" + file.Name()
			files = append(files, filepath)
		}
	}

}
func copyxml(path string) {
	file := strings.Replace(path, "----linux.log", ".xml", -1)
	file = strings.Replace(path, "采集完成目录", "xml", -1)
	ioutil.WriteFile(path, []byte(file), 0666)
}

func Runxml() {
	Dirs()
	if len(files) == 0 {
		log.Println("linux配置记录文件不存在哦。")
		return
	}
	for _, i := range files {
		ReplacepA1(i)
	}
}

// ReplacepA1 A1测评项
func ReplacepA1(file string) {
	xmlfilename := strings.Replace(file, "采集完成目录\\", "xml\\07.安全计算环境-服务器和终端-", -1)
	xmlfilename = strings.Replace(xmlfilename, "---linux.log", ".xml", -1)
	log.Println(xmlfilename)
	log.Println("-----")
	Passstr := "经核查，服务器在登录时采用用户名+口令的方式进行身份鉴别，"
	//错误数量，0为符合。4为不符合。其他为部分符合。
	errPASScoung := 0
	//整改建议。
	steer := ""
	//实际问题
	issue := ""
	//危害分析
	analyze := ""

	filedata, _ := ioutil.ReadFile(file)

	re := regexp.MustCompile(`不具备用户唯一性的用户为:(.*?)\n`)
	uidcount := re.FindAllStringSubmatch(string(filedata), -1)
	if len(uidcount[0][1]) == 0 {
		Passstr += "通过查看/etc/passwd账户文件中未发现同样的用户名和UID" + ";"
	} else {
		Passstr += "通过查看/etc/passwd账户文件中发现同样的用户名为:" + uidcount[0][1] + ";"
		errPASScoung += 1
		steer += "删除相同账户;"
		analyze += "存在相同用户可能导致在发生安全事件时无法确认责任人;"
	}

	re = regexp.MustCompile(`密码最大使用周期为:(.*?)\n`)
	passmax := re.FindAllStringSubmatch(string(filedata), -1)
	A1passmaxstr := passmax[0][1]
	intA1passmax, _ := strconv.Atoi(A1passmaxstr)
	switch {
	case intA1passmax > 90:
		errPASScoung += 1
		Passstr += "通过查看/etc/login.defs中PASS_MAX_DAYS为:" + A1passmaxstr + "定期更换密码周期过长;"
		steer += "设置密码更换周期为90天左右;"
		analyze += "账户口令可能会长期使用,恶意人员可通过猜测等方式获取口令;"
	case intA1passmax <= 90:
		Passstr += "通过查看/etc/login.defs中PASS_MAX_DAYS为:" + A1passmaxstr + "符合密码定期更换要求;"
	}
	re = regexp.MustCompile(`密码长度要求为:(.*?)\n`)
	passlen := re.FindAllStringSubmatch(string(filedata), -1)
	if len(passlen[0][1]) == 0 && len(passlen[0][1]) < 8 {
		errPASScoung += 1
		Passstr += "通过查看配置文件中未设置密码长度要求要求,不具备密码复杂度要求;"
		steer += "设置密码长度不小于为8位并设置密码复杂度要求。"
		//analyze += "账户口令可能被设置为弱口令恶意人员可通过暴力破解的方式获取账户口令，存在非授权访问的风险;"
	} else {
		ddluo := 0
		re = regexp.MustCompile(`密码数字位数要求为:(.*?)\n`)
		dcredit := re.FindAllStringSubmatch(string(filedata), -1)
		if len(dcredit[0][1]) != 0 {
			ddluo += 1
		}
		re = regexp.MustCompile(`密码小写位数要求为:(.*?)\n`)
		lcredit := re.FindAllStringSubmatch(string(filedata), -1)
		if len(lcredit[0][1]) != 0 {
			ddluo += 1
		}
		re = regexp.MustCompile(`密码大写位数要求为:(.*?)\n`)
		ucredit := re.FindAllStringSubmatch(string(filedata), -1)
		if len(ucredit[0][1]) != 0 {
			ddluo += 1
		}
		re = regexp.MustCompile(`密码特殊符号位数要求为:(.*?)\n`)
		ocredit := re.FindAllStringSubmatch(string(filedata), -1)
		if len(ocredit[0][1]) != 0 {
			ddluo += 1
		}
		if ddluo >= 2 {
			Passstr += "通过查看配置文件中设置密码长度要求为:" + passlen[0][1] + "密码数字位数要求为:" + dcredit[0][1] + "密码小写位数要求为:" + lcredit[0][1] + "密码大写位数要求为:" + ucredit[0][1] + "密码特殊符号位数要求为:" + ocredit[0][1] + "。"
		} else {
			Passstr += "通过查看配置文件中设置密码长度要求为:" + passlen[0][1] + "但为设置密码复杂度要求。"
			steer += "设置密码复杂度要求不少于为包含:数字、大小写字母以及特殊符号中的2种。"
			analyze += "账户口令可能被设置为弱口令恶意人员可通过暴力破解的方式获取账户口令，存在非授权访问的风险。"
		}
	}
	success := ""
	switch {
	case errPASScoung == 0:
		success += "符合"
	case errPASScoung == 4:
		success += "不符合"
	default:
		success += "部分符合"
	}
	//修改A1测评项
	issue = strings.Replace(steer, "设置", "未设置", -1)
	writedata = strings.Replace(gysxml, "A1核查结果", Passstr, -1)
	writedata = strings.Replace(writedata, "A1实际问题", issue, -1)
	writedata = strings.Replace(writedata, "A1符合情况", success, -1)
	writedata = strings.Replace(writedata, "A1整改建议", steer, -1)
	writedata = strings.Replace(writedata, "A1危害分析", analyze, -1)
	//修改A2测评项
	A2str, A2steer, A2issue, A2analyze, A2success := xml.A2(string(filedata))
	writedataA2 := strings.Replace(writedata, "A2核查结果", A2str, -1)
	writedataA2 = strings.Replace(writedataA2, "A2实际问题", A2issue, -1)
	writedataA2 = strings.Replace(writedataA2, "A2符合情况", A2success, -1)
	writedataA2 = strings.Replace(writedataA2, "A2整改建议", A2steer, -1)
	writedataA2 = strings.Replace(writedataA2, "A2危害分析", A2analyze, -1)
	//修改A3测评项
	A3str, A3steer, A3issue, A3analyze, A3success := xml.A3(string(filedata))
	writedataA3 := strings.Replace(writedataA2, "A3核查结果", A3str, -1)
	writedataA3 = strings.Replace(writedataA3, "A3实际问题", A3issue, -1)
	writedataA3 = strings.Replace(writedataA3, "A3符合情况", A3success, -1)
	writedataA3 = strings.Replace(writedataA3, "A3整改建议", A3steer, -1)
	writedataA3 = strings.Replace(writedataA3, "A3危害分析", A3analyze, -1)
	//修改B7测评项
	B7str, B7steer, B7issue, B7analyze, B7success := xml.B7(string(filedata))
	writedataB7 := strings.Replace(writedataA3, "B7核查结果", B7str, -1)
	writedataB7 = strings.Replace(writedataB7, "B7实际问题", B7issue, -1)
	writedataB7 = strings.Replace(writedataB7, "B7符合情况", B7success, -1)
	writedataB7 = strings.Replace(writedataB7, "B7整改建议", B7steer, -1)
	writedataB7 = strings.Replace(writedataB7, "B7危害分析", B7analyze, -1)
	//写入C3测评项
	C3str, C3steer, C3issue, C3analyze, C3success := xml.C3(string(filedata))
	writedataC3 := strings.Replace(writedataB7, "C3核查结果", C3str, -1)
	writedataC3 = strings.Replace(writedataC3, "C3实际问题", C3issue, -1)
	writedataC3 = strings.Replace(writedataC3, "C3符合情况", C3success, -1)
	writedataC3 = strings.Replace(writedataC3, "C3整改建议", C3steer, -1)
	writedataC3 = strings.Replace(writedataC3, "C3危害分析", C3analyze, -1)
	//写入D3测评项
	D3str, D3steer, D3issue, D3analyze, D3success := xml.D3(string(filedata))
	writedataD3 := strings.Replace(writedataC3, "D3核查结果", D3str, -1)
	writedataD3 = strings.Replace(writedataD3, "D3实际问题", D3issue, -1)
	writedataD3 = strings.Replace(writedataD3, "D3符合情况", D3success, -1)
	writedataD3 = strings.Replace(writedataD3, "D3整改建议", D3steer, -1)
	writedataD3 = strings.Replace(writedataD3, "D3危害分析", D3analyze, -1)
	//替换关联资产名称
	var xmlserver []string
	xmlserver = strings.Split(xmlfilename, "服务器和终端-")
	xmlserver = strings.Split(xmlserver[1], ".xml")
	log.Println(xmlserver[0])
	writedataD3 = strings.Replace(writedataD3, "高业尚服务器", xmlserver[0], -1)
	writefile(xmlfilename, writedataD3)

}

func writefile(path, data string) {
	_, err := os.Stat("xml")
	if os.IsNotExist(err) {
		os.Mkdir("xml", 0666)
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
	if err != nil {
		log.Println(err)
		return
	}
	f.Write([]byte(data))
}
