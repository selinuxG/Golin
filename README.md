
# 使用说明
> 此版本初衷是只做CLI版本,压缩体积小、跨平台可用，随着功能越多、参数越多使用人记不住参数，增加了基于Python的简易Gui版本，可满足基本增加资产、运行不同模式、运行Cli指令，增加了基于Web的多主机以及单主机不同模式的基本运行功能。注：完整使用功能建议使用Cli版本。


# 使用场景
> 自动化运维、多线程Linux、MySql、Redis、PostgreSQL、网络设备批量执行命令、等级保护（网络安全等级保护）现场测评工具、基线核查工具、测试syslog等

# 预览地址
> https://148.100.78.184:1818/golin/gys
## 子命令

### web
```shell
Usage:
  golin web [flags]

Flags:
  -h, --help          help for web
  -i, --ip string     指定运行网卡ip (default "127.0.0.1")
  -p, --port string   指定运行端口 (default "1818")
  -s, --save          是否额外保存文件
```

### update
```shell
通过api.github.com进行检查更新程序

Usage:
  golin update [flags]

Flags:
  -h, --help           help for update
  -p, --proxy string   此参数是指定代理ip(仅允许http/https代理哦)
```
### gui
```shell
通过python的tk开发,实现基本的增加资产以及运行,cli功能

Usage:
  golin gui [flags]

Flags:
  -h, --help   help for gui
```

### linux

```shell
基于SSH协议远程登陆,通过多线程的方法批量进行采集

Usage:
  golin linux [flags]

Flags:
  -c, --cmd string        此参数是指定待自定义执行的命令文件
  -C, --cmdvalue string   此参数是自定义执行命令（比-c优先级高）
  -e, --echo              此参数是控制控制台是否输出结果,默认不进行输出
  -h, --help              help for linux
  -i, --ip string         此参数是指定待远程采集的IP文件位置 (default "linux.txt")
  -l, --localhost         此参数是控制本机采集的模式
  -s, --spript string     此参数是指定IP文件中的分隔字符 (default "~")
  -v, --value string      此参数是指定执行单个主机
```

### mysql

```shell
基于Mysql远程通过多线程连接执行指定sql语句并记录。(也可本地执行)

Usage:
  golin mysql [flags]

Flags:
  -c, --cmd string      此参数是自定义执行sql语句
  -e, --echo            此参数指定是控制是否输出结果
  -h, --help            help for mysql
  -i, --ip string       此参数是指定待远程采集的IP文件位置 (default "mysql.txt")
  -s, --spript string   此参数是指定IP文件中的分隔字符 (default "~")
  -v, --value string    此参数是指定执行单个主机
```

### redis

```shell
基于Redis的远程登陆功能,通过多线程进行采集,基于info字段中的值判断,写入待采集文件主机时用户名为空即可。

Usage:
  golin redis [flags]

Flags:
  -h, --help            help for redis
  -i, --ip string       此参数是指定待远程采集的IP文件位置 (default "redis.txt")
  -s, --spript string   此参数是指定IP文件中的分隔字符 (default "~")
  -v, --value string    此参数是指定执行单个设备
```
### network
```shell
运行网络相关功能,目前仅有syslog模拟器

Usage:
  golin network [flags]

Flags:
  -h, --help     help for network
  -s, --syslog   模拟syslog接收端服务器

```
### route
```shell
基于SSH的功能进行采集

Usage:
  golin route [flags]

Flags:
  -c, --cmd string        此参数是指定待自定义执行的命令文件
  -C, --cmdvalue string   此参数是自定义执行命令（比-c优先级高）
  -e, --echo              此参数是控制控制台是否输出结果,默认不进行输出
  -h, --help              help for route
  -i, --ip string         此参数是指定待远程采集的IP文件位置 (default "route.txt")
  -p, --python            此参数是指定python位置，绝对路径，如'D:\python3\python.exe'
  -s, --spript string     此参数是指定IP文件中的分隔字符 (default "~")
```  
### 

### execl
```shell
通过读取xlsx文件生成golin可读取允许的格式文件

Usage:
  golin execl [flags]

Flags:
  -f, --file string     此参数是指定读取的文件
  -h, --help            help for execl
  -i, --ip string       此参数是指定ip代表的列
  -n, --name string     此参数是指定名称代表的列
  -p, --passwd string   此参数是指定密码所代表的列
  -P, --port string     此参数是指定端口代表的列
  -o, --sava string     此参数是指定保存的文件 (default "linux_xlsx.txt")
  -s, --sheet string    此参数是指定sheet名称 (default "Sheet1")
  -u, --user string     此参数是指定用户所代表的列
```
### pgsql
```
基于远程登录功能,通过多线程的方法批量进行采集(也可本地执行)
Usage:
  golin pgsql [flags]

Flags:
  -e, --echo            此参数是控制控制台是否输出结果,默认不进行输出
  -h, --help            help for pgsql
  -i, --ip string       此参数是指定待远程采集的IP文件位置 (default "pgsql.txt")
  -s, --spript string   此参数是指定IP文件中的分隔字符 (default "~")
  -v, --value string    此参数是单次执行

```
# 制作不同模式编辑启动脚本(bat)
## gui
```
@echo off
if "%1" == "h" goto begin
mshta vbscript:createobject("wscript.shell").run("%~nx0 h",0)(window.close)&&exit
:begin
golin.exe gui
```
## web
```
golin.exe web
```
