package poc

import (
	"golin/global"
	"strings"
)

var ListPocInfo []Flagcve

type Flagcve struct {
	Url  string
	Cve  string
	Flag string
}

func CheckPoc(url, app string) {

	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}

	functions := []func(string){
		AuthDruidun, //Druid未授权访问
		AuthSwagger, //swagger未授权访问
	}
	for _, function := range functions {
		function(url)
	}

	//以下是基于特征扫描
	app = strings.ToLower(app)
	switch {
	case strings.Contains(app, "spring"):
		CVE_2022_22947(url, "pwd")  //任意执行命令
		SpringActuatorHeapdump(url) //内存文件下载
	case strings.Contains(app, "nps"):
		NPS_default_passwd(url) //默认用户密码
	case strings.Contains(app, "struts2"):
		CVE_2021_31805(url) //任意执行命令
	}

}

func echoFlag(flag Flagcve) {
	global.PrintLock.Lock()
	defer global.PrintLock.Unlock()
	ListPocInfo = append(ListPocInfo, Flagcve{flag.Url, flag.Cve, flag.Flag})
}
