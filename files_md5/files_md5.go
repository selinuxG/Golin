package files_md5

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

var (
	oldpath     = flag.String("old", "nil", "指定原目录路径")
	newpath     = flag.String("new", "nil", "指定新目录路径")
	oldfiles    = make(map[string]string)
	newfiles    = make(map[string]string)
	oldfilesmd5 []string
	newfilesmd5 []string
)

func Run() {
	flag.Parse()
	log.Printf("------即将启动文件MD5对比功能，源目录为:%s,对比目录为:%s", *oldpath, *newpath)
	if *oldpath == "nil" {
		log.Println("源目录不能为空！！！")
		return
	}
	if *newpath == "nil" {
		log.Println("新目录路径不能为空！！！")
		return
	}
	GetFiles(*oldpath, "old")
	GetFiles(*newpath, "new")

	//检测MD5值
	for oldname, oldkey := range oldfiles {
		if chenkoldfilesmd5(oldkey) {
		} else {
			log.Printf("需要注意了！文件: %s 发生变化了！", oldname)
		}
	}

	//检测文件数量
	if len(newfiles) == len(oldfiles) {
	} else {
		if len(newfiles) > len(oldfiles) {
			log.Printf("需要注意了！文件数量多出%s个！", strconv.Itoa(len(newfiles)-len(oldfiles)))
			for _, v := range newfilesmd5 {
				if filecheck(v) == false {
					for i, k := range newfiles {
						if k == v {
							log.Println("多出文件:", i)
						}
					}
				}
			}
		}
		if len(oldfiles) > len(newfiles) {
			log.Printf("需要注意了！文件数量少出%s个！", strconv.Itoa(len(oldfiles)-len(newfiles)))

		}
	}
	//监测新文件与
}

//切片读取多出的文件对比
func filecheck(data string) bool {
	for _, v := range newfiles {
		if v == data {
			return true
		}
	}
	return false
}

//递归读取目录下文件
func GetFiles(folder string, par string) {
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			GetFiles(folder+"/"+file.Name(), par)
		} else {
			firename := folder + "/" + file.Name()
			md5flie, err := FileMD5(firename)
			if err != nil {
				log.Println(err)
			}
			if par == "new" {
				newfiles[firename] = md5flie
				newfilesmd5 = append(newfilesmd5, md5flie)
			}
			if par == "old" {
				oldfiles[firename] = md5flie
				//oldfiles = append(oldfiles, firename)
				oldfilesmd5 = append(oldfilesmd5, md5flie)
			}
		}
	}
}

func FileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	_, _ = io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil)), nil
}

//检查源md5是否存在
func chenkoldfilesmd5(md string) bool {
	for _, v := range newfiles {
		if md == v {
			return true
		}
	}
	return false
}
