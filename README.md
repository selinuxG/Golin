# CLI版

> 目前在从写阶段，V1版本主函数将近2000行代码，很难维护，并且存在诸多不足以及存在重复代码，此版本目的就是为了解决这个问题，并优化使用逻辑，做成一个正儿八经的命令行工具，并删除无用功能。
正在疯狂敲打键盘中....

## 子命令

### linux

```shell
./golin linux -h

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
./golin mysql -h

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
./golin redis -h

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
基于google发布的库进行抓取流量
过滤源 IP 地址为 192.168.1.1 的数据包：src host 192.168.1.1  过滤目标 IP 地址为 192.168.1.1 的数据包：dst host 192.168.1.1
过滤源或目标 IP 地址为 192.168.1.1 的数据包：host 192.168.1.1  过滤源端口为 80 的数据包：src port 80
过滤目标端口为 80 的数据包：dst port 80  过滤 TCP 数据包：tcp  过滤 UDP 数据包：udp  组合 多个条件进行过滤：tcp and dst port 80 and src host 192.168.1.1

Usage:
  golin network [flags]

Flags:
  -b, --bfp string         指定监听网卡 (default "tcp")
  -h, --help               help for network
  -i, --interface string   指定监听网卡 (default "\\Device\\NPF_{B3C9B1B3-FFC0-4D59-8667-4B0BF6D354CE}")
  -a, --interfaceall       显示本地所有网卡及IP地址
```
### 

# V1版

## 工具介绍

> 在做等保以及其他一部分工作时，经常会遇到几十甚至上百台服务器。
人工一个一个尝试登陆在执行命令保存结果时间太长，于此是此工具一部分功能的由来。
执行默认为等保测评项的所需命令，并保存输出结果。可通过参数使用自定义命令，只要设备为SSH协议登录都正常用。
我希望这个小巧的工具具有自动化测评的能力以及具备一部分安全的功能。


## 等保功能

1. -run linux、redis、mysql、postgresql、oracle、windos	"自动化扫描类型主机"
2. -run oracle  -con system/oracle@1.1.2.135:1521/sid -name oracle数据库 “oracle为特定参数”
3. -cmd cmd.txt	"自定义执行指定文件中的命令，多个命令";"分隔。
4. -db linux	"输出等保测评的常见设备命令，现支持oracle、aix、huawei、mysql、linux、达梦、cisco、postgresql、nginx、mongo"
5. -checkpass true "需要配合-run参数使用，验证密码是否具备复杂度要求(字⺟⼤⼩写+数字+特殊符号，8位以上)"

## 安全检测功能

1. -webserver www.baidu.com “输出指定地址的开放的web服务并html保存源码，如是windos则弹窗源码目录”
2. -webserver www.baidu.com -img true “输出指定地址的开放的web服务并截图，如是windos则弹窗图片目录”
3. -filesmd5 true  -old D:/testbak -new D:/test	"对比两个目录中文件的MD5值，确认文件是否被更改"
4. -port  www.baidu.com “输出指定域名的开放端口”
5. -ifconfig true  “输出当前外网地址”
6. -ipinfo 47.94.159.75   “输出ip地址信息，所在地，经纬度”
7. -systeminfo true  “输出系统信息”

## 其他功能参数

1. -fileshare true  “通过web方式共享当前目录,true开启,端口为11111”
