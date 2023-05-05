package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golin/config"
	"golin/global"
	"os/exec"
	"path/filepath"
)

// guiCmd represents the gui command
var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "运行简易版GUI辅助程序",
	Long:  `通过python的tk开发,实现基本的增加资产并允许功能`,
	Run: func(cmd *cobra.Command, args []string) {
		guipath := filepath.Join(global.PythonDir, global.PyGui)
		if !global.PathExists(guipath) {
			config.Log.Warn("gui文件不存在!", zap.String("文件:", guipath))
			return
		}
		Python_path := global.PythonPath
		// 接收第一个传参为python执行目录
		if len(args) >= 1 {
			if !global.PathExists(args[0]) {
				config.Log.Warn("python程序不存在!", zap.String("python-path:", args[0]))
				return
			}
			Python_path = args[0]
		}

		runcmd := exec.Command(Python_path, guipath)
		err := runcmd.Run()
		if err != nil {
			config.Log.Warn("运行Gui.py失败!")
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(guiCmd)
}
