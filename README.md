# 使用场景
> 自动化运维、批量执行命令、等保（网络安全等级保护）测评工具、基线核查工具
# CLI版

> 此版本主做cli，目的是为了跨平台可用，并附带简易版Gui程序供选择。

## 子命令
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
基于Mysql远程通过多线程连接执行指定sql语句并记录。

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