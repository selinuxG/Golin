# 工具介绍
> 在做等保以及其他一部分工作时，经常会遇到几十甚至上百台服务器。
> 人工一个一个尝试登陆在执行命令保存结果时间太长，于此是此工具一部分功能的由来。
> 执行默认为等保测评项的所需命令，并保存输出结果。可通过参数使用自定义命令，只要设备为SSH协议登录都正常用。
> 我希望这个小巧的工具具有自动化测评的能力以及具备一部分安全的功能。

# 等保功能

1. -run linux、redis、mysql、postgresql、oracle、windos	"自动化扫描类型主机"
1. -run oracle  -con system/oracle@1.1.2.135:1521/sid -name oracle数据库 “oracle为特定参数”
1. -cmd cmd.txt	"自定义执行指定文件中的命令，多个命令";"分隔。
1. -db linux	"输出等保测评的常见设备命令，现支持oracle、aix、huawei、mysql、linux、达梦、cisco、postgresql、nginx、mongo"
1. -checkpass true "需要配合-run参数使用，验证密码是否具备复杂度要求(字⺟⼤⼩写+数字+特殊符号，8位以上)"
1. 由于官网未提供SDK自动化采集oracle存在不稳定性，经本地测评Oracle11无问题。
# 安全检测功能

1. -webserver www.baidu.com “输出指定地址的开放的web服务并html保存源码，如是windos则弹窗源码目录”
1. -webserver www.baidu.com -img true “输出指定地址的开放的web服务并截图，如是windos则弹窗图片目录”
1. -filesmd5 true  -old D:/testbak -new D:/test	"对比两个目录中文件的MD5值，确认文件是否被更改"
1. -port  www.baidu.com “输出指定域名的开放端口”
1. -ifconfig true  “输出当前外网地址”
1. -ipinfo 47.94.159.75   “输出ip地址信息，所在地，经纬度”
1. -systeminfo true  “输出系统信息”
# 其他功能参数

1. -fileshare true  “通过web方式共享当前目录,true开启,端口为11111”
1. -imgbing true “获取今日bing官网壁纸”
# 使用方法
> 扫描任务来源于当面目录下的：ip.txt(linux)，mysql.txt,redis.txt,postgresql.txt。
> 一行一个扫描任务，格式为:(空密码或空密码的话为空即可)
>
> 服务器名称~ip~user~password~port
> 默认为并发扫描，测试100+linux服务器大约用时5秒。
> 增加扫描信息可配合另一个GUI工具自动添加。

## windos
> cmd下运行golin.exe文件，通过参数控制运行。

## Linux
> 先赋予执行权限，在执行。
> chmod 777 golin (权限自己控制，确保具有当前目录创建以及执行权限)
> 运行
> ./golin -run linux(mysql、redis、postgresql)
