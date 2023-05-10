package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"golin/global"
	"golin/run"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

var save bool

func Start(cmd *cobra.Command, args []string) {
	ip, _ := cmd.Flags().GetString("ip")
	port, _ := cmd.Flags().GetString("port")
	save, _ = cmd.Flags().GetBool("save")
	r := gin.Default()
	r.SetHTMLTemplate(IndexTemplate()) //加载首页
	golin := r.Group("/golin")
	{
		golin.GET("/index", GolinIndex)    //首页
		golin.POST("/submit", GolinSubmit) //提交任务
	}
	r.Run(ip + ":" + port)
}

// GolinIndex 首页
func GolinIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index", "")
}

// GolinSubmit 提交任务
func GolinSubmit(c *gin.Context) {
	name, ip, user, passwd, port, mode, down := c.PostForm("name"), c.PostForm("ip"), c.PostForm("user"), c.PostForm("password"), c.PostForm("port"), c.PostForm("run_mode"), c.PostForm("down")
	//fmt.Println(name, ip, user, passwd, port, mode)
	savefilename := fmt.Sprintf("%s_%s.log", name, ip)                //保存的文件夹名：名称_ip.log
	successfile := filepath.Join(global.Succpath, mode, savefilename) //保存的完整路径
	if global.PathExists(successfile) {
		c.String(http.StatusOK, "保存的文件中有重名文件，更换一个吧客官~")
		c.Abort()
		return
	}
	run.Onlyonerun(fmt.Sprintf("%s~~~%s~~~%s~~~%s~~~%s", name, ip, user, passwd, port), "~~~", mode)
	if global.PathExists(successfile) {
		//如果不保存文件，文件返回后删除
		defer func() {
			if !save {
				os.Remove(successfile)
			}
		}()
		//down下载文件,preview预览文件
		if down == "down" {
			c.Header("Content-Description", "File Transfer")
			c.Header("Content-Disposition", "attachment; filename="+fmt.Sprintf(fmt.Sprintf("%s_%s(%s).log", name, ip, mode)))
			c.Header("Content-Type", "application/octet-stream")
		}
		//返回文件
		c.File(successfile)
	} else {
		c.String(200, "失败了哦客官～")
	}

}

// IndexTemplate 返回包含模板内容的模板结构体
func IndexTemplate() *template.Template {
	tmpl, err := template.New("index").Parse(IndexHtml())
	if err != nil {
		panic(err)
	}
	return tmpl
}
