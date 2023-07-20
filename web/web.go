package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"golin/global"
	"os/exec"
	"runtime"
)

var save bool

func Start(cmd *cobra.Command, args []string) {

	if !global.PathExists("cert.pem") || !global.PathExists("cert.key") {
		CreateCert()
	}

	ip, _ := cmd.Flags().GetString("ip")
	port, _ := cmd.Flags().GetString("port")
	save, _ = cmd.Flags().GetBool("save")
	r := gin.Default()
	r.NoRoute(func(c *gin.Context) {
		GolinErrorhtml("404", "sorry~请求不存在哦!", c)
	})
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/golin/gys")
	})

	golin := r.Group("/golin")
	{
		golin.GET("/gys", GolinHome)               //首页
		golin.GET("/index", GolinIndex)            //单主机index
		golin.GET("/indexfile", GolinIndexFile)    //多主机index
		golin.GET("/modefile", GolinMondeFileGet)  //返回模板文件
		golin.POST("/submit", GolinSubmit)         //提交单主机任务
		golin.POST("/submitfile", GolinSubmitFile) //提交多主机任务
		golin.GET("/history", GolinHistory)        //历史记录
		golin.GET("/update", GolinUpdate)          //检查更新
		golin.GET("/dj", GolinDj)                  //模拟定级首页
		golin.POST("/djPost", GolinDjPost)         //提交定级结果

	}
	// Windows下在默认浏览器中打开网页
	go func() {
		if runtime.GOOS == "windows" {
			cmd := exec.Command("cmd", "/C", fmt.Sprintf("start https://%s:%s/golin/gys", ip, port))
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error opening the browser:", err)
			}
		}
	}()
	// 启动gin
	r.RunTLS(ip+":"+port, "cert.pem", "key.pem")
}
