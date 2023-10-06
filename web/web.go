package web

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"golin/global"
	"net/http"
	"os/exec"
	"runtime"
)

var save bool

//go:embed favicon.ico
var faviconFS embed.FS

func Start(cmd *cobra.Command, args []string) {

	if !global.PathExists("cert/cert.pem") || !global.PathExists("cert/cert.key") {
		CreateCert()
	}

	ip, _ := cmd.Flags().GetString("ip")
	port, _ := cmd.Flags().GetString("port")
	save, _ = cmd.Flags().GetBool("save")
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Header("author", "gaoyeshang")
		c.Header("VX", "SelinuxG")
		c.Next()
	})
	r.NoRoute(func(c *gin.Context) {
		GolinErrorhtml("404", "sorry~请求不存在哦!", c)
	})
	r.GET("/favicon.ico", faviconHandler) //路由图标
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
	// Windows、Mac下在默认浏览器中打开网页
	go func() {
		if runtime.GOOS == "windows" {
			cmd := exec.Command("cmd", "/C", fmt.Sprintf("start https://%s:%s/golin/gys", ip, port))
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error opening the browser:", err)
			}
		}
		if runtime.GOOS == "darwin" {
			cmd := exec.Command("open", fmt.Sprintf("https://%s:%s/golin/gys", ip, port))
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error opening the browser:", err)
			}
		}
	}()
	// 启动gin
	r.RunTLS(ip+":"+port, "cert/cert.pem", "cert/key.pem")
}

func faviconHandler(c *gin.Context) {
	c.FileFromFS("favicon.ico", http.FS(faviconFS))
}
