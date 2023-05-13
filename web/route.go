package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golin/global"
	"golin/run"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// GolinIndex 单主机首页
func GolinIndex(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, IndexHtml())
}

// GolinIndexFile 多主机首页
func GolinIndexFile(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, IndexFilehtml())
}

func GolinSubmitFile(c *gin.Context) {
	filename, err := c.FormFile("uploaded-file")
	if err != nil {
		c.String(http.StatusOK, "上传文件失败！")
		c.Abort()
		return
	}
	tempfilenamexlsx := fmt.Sprintf("%v.xlsx", time.Now().Unix())
	tempfilenametxt := fmt.Sprintf("%v.txt", time.Now().Unix())
	defer os.Remove(tempfilenamexlsx)
	defer os.Remove(tempfilenametxt)

	err = c.SaveUploadedFile(filename, tempfilenamexlsx)
	if err != nil {
		c.String(http.StatusOK, "上传保存失败！")
		c.Abort()
		return
	}
	if CreateTmpTxt(tempfilenamexlsx, tempfilenametxt) {
		mode := c.PostForm("mode")
		run.Rangefile(tempfilenametxt, "~~", mode) //运行多主机采集模式
		if global.PathExists(filepath.Join(global.Succpath, mode)) {
			if runtime.GOOS == "windows" {
				ip := c.RemoteIP()
				iplist, _ := global.GetLocalIPAddresses()
				for _, v := range iplist {
					if ip == v || ip == "127.0.0.1" {
						cmd := exec.Command("explorer.exe", filepath.Join(global.Succpath, mode))
						cmd.Run()
						break
					}
				}
			}
		} else {
			c.String(http.StatusOK, "运行失败了哦！")
			c.Abort()
			return
		}
		c.Redirect(302, "/golin/indexfile")
	}
}

// GolinSubmit 单词提交任务
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

	switch mode {
	case "Route": //路由模式是单独的
		for _, cmd := range run.Defroutecmd {
			run.Routessh(successfile, ip, user, passwd, port, cmd)
		}
	default: //其他模式统一函数传参
		run.Onlyonerun(fmt.Sprintf("%s~~~%s~~~%s~~~%s~~~%s", name, ip, user, passwd, port), "~~~", mode)
	}

	//run.Onlyonerun(fmt.Sprintf("%s~~~%s~~~%s~~~%s~~~%s", name, ip, user, passwd, port), "~~~", mode)
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

// GolinMondeFileGet 返回模板文件
func GolinMondeFileGet(c *gin.Context) {
	//如果本地没有模板文件则生成一个
	if !global.PathExists(global.XlsxTemplateName) && !CreateTemplateXlsx() {
		c.String(http.StatusOK, "模板文件生成失败！")
	}
	// 添加必要的响应头以触发文件下载
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+global.XlsxTemplateName)
	c.Header("Content-Type", "application/octet-stream")
	c.File(global.XlsxTemplateName)
}
