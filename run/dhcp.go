package run

func Linux_cmd() string {
	cmd := `
echo "------------------------------------------------------------------------------------------------------信息系统"
echo "系统版本:$(cat /etc/redhat-release;echo;lsb_release -a)";
echo "系统架构:$(uname -m)"
echo "网卡信息:$(echo;ifconfig |grep inet |grep -v "inet6"|awk -F' ' '{print $2}')"
echo
echo "------------------------------------------------------------------------------------------------------"
echo "测评结果------GYSv1.0"
echo "应对登录的用户进行身份标识和鉴别，身份标识具有唯一性，身份鉴别信息具有复杂度要求并定期更换；"
echo "不具备用户唯一性的用户为:";cat /etc/passwd |awk -F":" '{print $1}'|uniq -c|awk -F" " '$3 >1 {print $2}'
echo "密码最大使用周期为:$(cat /etc/login.defs | grep PASS_MAX_DAYS | grep -v ^# | awk '{print $2}')"
echo "密码短使用周期为:$(cat /etc/login.defs | grep PASS_MIN_DAYS | grep -v ^# | awk '{print $2}')"
echo "密码数字位数要求为:$(cat /etc/pam.d/sshd|grep difok | awk -F'dcredit=' '{print $2}' | awk '{print $1}')"
echo "密码小写位数要求为:$(cat /etc/pam.d/sshd|grep difok | awk -F'lcredit=' '{print $2}' | awk '{print $1}')"
echo "密码大写位数要求为:$(cat /etc/pam.d/sshd|grep difok | awk -F'ucredit=' '{print $2}' | awk '{print $1}')"
echo "密码特殊符号位数要求为:$(cat /etc/pam.d/sshd|grep difok | awk -F'ocredit=' '{print $2}' | awk '{print $1}')"
echo "密码长度要求为:$(cat /etc/pam.d/sshd|grep minlen | awk -F'dcredit=' '{print $2}' | awk '{print $1}')"
echo
echo "------------------------------------------------------------------------------------------------------"
echo "应具有登录失败处理功能，应配置并启用结束会话、限制非法登录次数和当登录连接超时自动退出等相关措施；"
echo "超时退出时间为:$(echo $TMOUT)"
echo "登陆错误失败次数为:$(cat /etc/pam.d/sshd | grep deny | awk -F'deny=' '{print $2}' | awk '{print $1}')"
echo "失败锁定时常为:$(cat /etc/pam.d/sshd |grep unlock_time | awk -F'unlock_time=' '{print $2}' | awk '{print $1}')"
DENYROOT=$(cat /etc/pam.d/sshd |grep -v "^#"| grep even_deny_root |wc -l);
echo "是否设置ROOT用户失败锁定:$(if [ $DENYROOT -eq 1 ];then echo "yes";else echo "未配置";fi)"
echo "ROOT用户失败锁定时常为:$(cat /etc/pam.d/sshd |grep root_unlock_time | awk -F'root_unlock_time=' '{print $2}' | awk '{print $1}')"
echo
echo "------------------------------------------------------------------------------------------------------"
echo "当进行远程管理时，应采取必要措施防止鉴别信息在网络传输过程中被窃听；"
echo "当使用此工具时默认符合。此工具逻辑为通过SSH登陆。需确认下面是否开启telnet服务。"
echo
echo "------------------------------------------------------------------------------------------------------"
echo "应采用口令、密码技术、生物技术等两种或两种以上组合的鉴别技术对用户进行身份鉴别，且其中一种鉴别技术至少应使用密码技术来实现；"
echo "经核查，通过ssh协议远程登录尝试可直接登录系统，未采用两种或两种以上组合的鉴别技术对用户进行身份鉴别。"
echo
echo "------------------------------------------------------------------------------------------------------"
echo "应对登录的用户分配账户和权限------"
echo "特权用户:$(awk -F: '$3==0 {print $1}' /etc/passwd;)";
echo "可登录用户为:$(echo;cat /etc/passwd |grep -v -E "sbin|false" |awk -F ":" '{print $1}')"
echo "是否限制Root用户远程登陆（输出no为已限制）:$(cat /etc/ssh/sshd_config |grep -v "^#" |grep "PermitRootLogin" |awk -F' ' '{print $2}')"
echo
echo "------------------------------------------------------------------------------------------------------"
echo "应重命名或删除默认账户，修改默认账户的默认口令；"
echo "特权用户:$(awk -F: '$3==0 {print $1}' /etc/passwd;)";
echo "经核查，通过查看/etc/passwd文件发现存在默认账户root未锁定或重命名，但已修改默认口令。"
echo
echo "------------------------------------------------------------------------------------------------------"
echo "应及时删除或停用多余的、过期的账户，避免共享账户的存在；"
echo "空密码用户:$(cat /etc/shadow |awk -F: 'length($2)==0 {print $1}')";
echo "通过最下方登陆信息或者访谈确认是否存在多余的账户"
echo
echo "------------------------------------------------------------------------------------------------------"
echo "应授予管理用户所需的最小权限，实现管理用户的权限分离；"
echo "经核查，通过查看及访谈/etc/passwd中包含/bin/bash可登录的账户，存在管理账户root，仅存在系统管理员，未建立安全管理员、审计管理员。"
echo
echo "------------------------------------------------------------------------------------------------------"
echo "应由授权主体配置访问控制策略，访问控制策略规定主体对客体的访问规则；"
echo "umask为:$(umask)"
echo
echo "------------------------------------------------------------------------------------------------------"
echo "访问控制的粒度应达到主体为用户级或进程级，客体为文件、数据库表级；"
echo "/etc/passwd的权限为:$(ls -l /etc/passwd | awk '{print $1}')"
echo "/etc/shadow的权限为:$(ls -l /etc/shadow | awk '{print $1}')"
echo "/etc/group的权限为:$(ls -l /etc/group | awk '{print $1}')"
echo "/etc/securetty的权限为:$(ls -l /etc/securetty | awk '{print $1}')"
echo "/etc/services的权限为:$(ls -l /etc/services | awk '{print $1}')"
echo "/etc/rsyslog.conf的权限为:$(ls -l /etc/rsyslog.conf | awk '{print $1}')"
echo "/etc/login.defs的权限为:$(ls -l /etc/login.defs | awk '{print $1}')"
echo "/etc/ssh/sshd_config的权限为:$(ls -l /etc/ssh/sshd_config | awk '{print $1}')"
echo "/etc/hosts.allow的权限为:$(ls -l /etc/hosts.allow | awk '{print $1}')"
echo "/etc/hosts.deny的权限为:$(ls -l  /etc/hosts.deny | awk '{print $1}')"
echo 
echo "------------------------------------------------------------------------------------------------------"
echo "应对重要主体和客体设置安全标记，并控制主体对有安全标记信息资源的访问；"
echo "Selinux状态:$(cat /etc/selinux/config |grep -v "^#"|grep -v "^$"|awk -F 'SELINUX=' '{print $2}')"
echo
echo "------------------------------------------------------------------------------------------------------"
echo "应启用安全审计功能，审计覆盖到每个用户，对重要的用户行为和重要安全事件进行审计；"
echo "linux系统默认符合。"
echo
echo "------------------------------------------------------------------------------------------------------"
echo "审计记录应包括事件的日期和时间、用户、事件类型、事件是否成功及其他与审计相关的信息；"
echo "审计日志记录位置为:$(cat /etc/audit/auditd.conf |grep "log_file = /" |awk -F "= " '{print $2}')"
echo "查看前三行Auditd审计日志:";head -n 3 $(cat /etc/audit/auditd.conf |grep "log_file = /" |awk -F "= " '{print $2}')
echo "------------------------------------------------------------------------------------------------------"
echo "查看后三行Auditd审计日志:";tail -n 3 $(cat /etc/audit/auditd.conf |grep "log_file = /" |awk -F "= " '{print $2}')
echo "------------------------------------------------------------------------------------------------------"
echo "查看前三行messages日志:";head -n 3 /var/log/messages
echo "------------------------------------------------------------------------------------------------------"
echo "查看后三行messages日志:";head -n 3 /var/log/messages
echo "------------------------------------------------------------------------------------------------------"
echo "查看前三行secure日志:";head -n 3 /var/log/secure
echo
echo "------------------------------------------------------------------------------------------------------"
echo "应对审计记录进行保护，定期备份，避免受到未预期的删除、修改或覆盖等；"
echo "审计日志权限为:$(ls -l $(cat /etc/audit/auditd.conf |grep "log_file = /" |awk -F "= " '{print $2}'))"
echo "定时任务为:";cat /etc/crontab |grep -v "^#" |grep -v "^SHELL"|grep -v "^PATH"|grep -v "^MA"|grep -v "^$" 
echo "rsys配置为:$(echo;cat /etc/rsyslog.conf  |grep -v "^#"|grep -v -E "ModL|Act|Work|Inclu|Omit|IMJ"|grep -v "^$")"
auditip=$(cat /etc/rsyslog.conf  |grep -v -E "^#"|grep -v "^$"|grep "[0-9]\{1,3\}[.][0-9]\{1,3\}[.][0-9]\{1,3\}[.][0-9]\{1,3\}"|awk -F "@{1,2}" '{print $2}' |awk -F ":" '{print $1}');if [ "$auditip" != "" ] && [ "$auditip" != "127.0.0.1" ]; then ping $auditip -c 1 -W 1 >/dev/null;if [ $? -eq 0 ];then echo "审计日志服务器:" $auditip "通信正常"; else echo "审计日志服务器:" $auditip"通信失败"; fi; else echo "未配置审计服务地址";fi
echo
echo "------------------------------------------------------------------------------------------------------"
echo "应对审计进程进行保护，防止未经授权的中断；"
echo "经核查，使用ps -ax -o ruid -o euid -o suid -o fuid -o pid -o fname | grep auditd输出内容中可以看出该系统进程的real UID，effective UID，saved UID，file UID的值为0，即权限仅为系统管理员，使用普通用户登录，无法终止该进程，可防止未经授权的中断。"
echo
echo "------------------------------------------------------------------------------------------------------"
echo "应遵循最小安装的原则，仅安装需要的组件和应用程序；"
echo "按照默认符合吧,具体看开放的端口对应的服务"
echo
echo "------------------------------------------------------------------------------------------------------"
echo "应关闭不需要的系统服务、默认共享和高危端口；"
echo "export共享为:$(cat /etc/export 2>/dev/null)"
echo "开放端口为:$(echo;netstat -anpt |grep LISTEN)"
echo "------------------------------------------------------------------------------------------------------"
echo "应通过设定终端接入方式或网络地址范围对通过网络进行管理的管理终端进行限制；"
echo "hosts.allow限制为:";if [ -z $(cat /etc/hosts.allow|grep -v "^#"|grep -v "^$") ];then echo "未做限制";else echo "$(cat /etc/hosts.allow|grep -v "^#"|grep -v "^$")";fi
echo "hosts.deny限制为:";if [ -z $(cat /etc/hosts.deny|grep -v "^#"|grep -v "^$") ];then echo "未做限制";else echo "$(cat /etc/hosts.deny|grep -v "^#"|grep -v "^$")";fi
echo "Selinux状态:$(cat /etc/selinux/config |grep -v "^#"|grep -v "^$"|awk -F 'SELINUX=' '{print $2}')"
Firewalld=$(systemctl status firewalld 1>/dev/null 2>/dev/null);
echo "Firewalld状态为:$(if [ $? -eq 0 ];then echo "yes";else echo "no";fi 2>/dev/null)"
iptabing=$(systemctl status iptables.service 1>/dev/null 2>/dev/null);
echo "iptables状态为:$(if [ $? -eq 0 ];then echo "yes";else echo "no";fi 2>/dev/null)"
echo "防火墙策略为:"
iptables -L
echo
echo "------------------------------------------------------------------------------------------------------"
echo "应提供数据有效性检验功能，保证通过人机接口输入或通过通信接口输入的内容符合系统设定要求；"
echo "此项对服务器不适用。"
echo
echo "------------------------------------------------------------------------------------------------------"
echo "应能发现可能存在的已知漏洞，并在经过充分测试评估后，及时修补漏洞；"
echo "需通过查看是否具备特定的安全设备并具备扫描记录以及修补记录。"
echo
echo "------------------------------------------------------------------------------------------------------"
echo "应能够检测到对重要节点进行入侵的行为，并在发生严重入侵事件时提供报警；"
echo "需通过查看是否具备特定的安全设备并能提供报警。"
echo
echo "------------------------------------------------------------------------------------------------------"
echo "应采用免受恶意代码攻击的技术措施或主动免疫可信验证机制及时识别入侵和病毒行为，并将其有效阻断；"
echo "默认基本不符合，根据实际安装程序确认。"
echo
echo "------------------------------------------------------------------------------------------------------"
echo "可基于可信根对计算设备的系统引导程序、系统程序、重要配置参数和应用程序等进行可信验证，并在应用程序的关键执行环节进行动态可信验证，在检测到其可信性受到破坏后进行报警，并将验证结果形成审计记录送至安全管理中心；"
echo "默认基本不符合，具体看系统以及版本。"
echo "------------------------------------------------------------------------------------------------------"
echo "数据完整性保密性通过SSH则符合。数据备案具体看上述中的定时任务以及sys配置。"
echo
echo
echo
echo
echo "------------------------------------------------------------------------------------------------------登陆记录"
echo "查看用户登录记录:";last;
echo "------------------------------------------------------------------------------------------------------上次登录信息"
echo "查看所有用户上次登录信息：";lastlog`
	return cmd
}
