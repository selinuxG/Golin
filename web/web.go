package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"html/template"
	"os/exec"
	"runtime"
)

var save bool

func Start(cmd *cobra.Command, args []string) {
	ip, _ := cmd.Flags().GetString("ip")
	port, _ := cmd.Flags().GetString("port")
	save, _ = cmd.Flags().GetBool("save")
	r := gin.Default()
	r.NoRoute(func(c *gin.Context) {
		GolinErrorhtml("404", "sorry~请求不存在哦!", c)
	})
	golin := r.Group("/golin")
	{
		golin.GET("/index", GolinIndex)            //首页
		golin.GET("/indexfile", GolinIndexFile)    //多主机首页
		golin.GET("/modefile", GolinMondeFileGet)  //返回模板文件
		golin.POST("/submit", GolinSubmit)         //提交单主机任务
		golin.POST("/submitfile", GolinSubmitFile) //提交多主机任务
	}
	// Windows下在默认浏览器中打开网页
	go func() {
		if runtime.GOOS == "windows" {
			cmd := exec.Command("cmd", "/C", fmt.Sprintf("start http://%s:%s/golin/indexfile", ip, port))
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error opening the browser:", err)
			}
		}
	}()
	// 启动gin
	r.Run(ip + ":" + port)
}

// Template 返回包含模板内容的模板结构体
func Template(html, indexname string) *template.Template {
	tmpl, err := template.New(indexname).Parse(html)
	if err != nil {
		panic(err)
	}
	tmpl = tmpl.Lookup(indexname)
	return tmpl
}
