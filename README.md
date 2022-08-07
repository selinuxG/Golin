# 工具介绍

```
在做等保以及其他一部分工作时，经常会遇到几十甚至上百台服务器，人工一个一个尝试登陆在执行命令保存结果时间太长，于此是此工具的由来。执行默认为等保测评项的所需命令，并保存输出结果。
支持Linux、Mysql、Redis、Postgresql。
```

# 使用方法

## windos

```cmd
cmd下运行golin.exe文件，通过参数控制运行。
golinx.exe -run linux(mysql、redis、postgresql)
```

## Linux

```bash
先赋予执行权限，在执行。
chmod 777 golin (权限自己控制，确保具有当前目录创建以及执行权限)
运行
./golin -run linux(mysql、redis、postgresql)
```



# 等保自动化采集参数

​	通过cmd或者bash执行

- -run linux “采集linux服务器”
- -run linux -cmd cmd.txt  “通过指定命令文件采集允许ssh的服务器”（注：命令为一行通过；分隔）
- -run mysql “采集mysql数据库”
- -run redis “采集redis数据库”
- -run postgresql “采集postgresql数据库”
- -db huawei  “输出等保华为的命令，支持linux，达梦，oracle，mysql，cisco，huawei，aix，postgresql”

 #  其他功能参数

- -webserver www.baidu.com “输出指定地址的开放的web服务并html保存源码，如是windos则弹窗源码目录”
- -webserver www.baidu.com -img true “输出指定地址的开放的web服务并截图，如是windos则弹窗图片目录”
- -port  www.baidu.com “输出指定域名的开放端口”
- -ifconfig true  “输出当前外网地址”
- -ipinfo 47.94.159.75   “输出ip地址信息，所在地，经纬度”
- -systeminfo true  “输出系统信息”
- -fileshare true  “通过web方式共享当前目录,true开启,端口为11111”
- -imgbing true “获取今日bing官网壁纸”

