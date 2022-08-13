package dbcp

func CMD1() string {
	cmd := `
echo "------------------------------------------------------------------------------------------------------版本信息"
echo "系统版本:$(cat /etc/redhat-release)";
echo "系统架构:$(uname -m)"
echo "网卡信息:$(echo;ifconfig |grep inet |grep -v "inet6"|awk -F' ' '{print $2}')"
echo "------------------------------------------------------------------------------------------------------检查服务"
echo "$(if [ $(ps aux |grep nginx |grep -v "color"|wc -l) -gt 0 ];then echo "存在nginx服务";fi)"
echo "$(if [ $(ps aux |grep mysql |grep -v "color"|wc -l) -gt 0 ];then echo "存在mysql服务";fi)"
echo "$(if [ $(ps aux |grep redis |grep -v "color"|wc -l) -gt 0 ];then echo "存在redis服务";fi)"
echo "$(if [ $(netstat -anpt |grep java |grep -v "color"|wc -l) -gt 0 ];then echo "正在运行java服务存在tomcat服务";fi)"
echo "$(if [ $(netstat -anpt |grep mongod |grep -v "color"|wc -l) -gt 0 ];then echo "存在mongod服务";fi)" 
echo "$(if [ $(netstat -anpt |grep ora |grep -v "color"|wc -l) -gt 0 ];then echo "存在oracle服务";fi)"


echo "------------------------------------------------------------------------------------------------------用户信息"
echo "不具备用户唯一性的用户为:";cat /etc/passwd |awk -F":" '{print $1}'|uniq -c|awk -F" " '$3 >1 {print $2}'
echo "特权用户:$(awk -F: '$3==0 {print $1}' /etc/passwd;)";
echo "空密码用户:$(cat /etc/shadow |awk -F: 'length($2)==0 {print $1}')";
echo "可登录用户为:$(echo;cat /etc/passwd |grep -v -E "sbin|false" |awk -F ":" '{print $1}')"
echo "密码最大使用周期为:$(cat /etc/login.defs | grep PASS_MAX_DAYS | grep -v ^# | awk '{print $2}')";
echo "密码短使用周期为:$(cat /etc/login.defs | grep PASS_MIN_DAYS | grep -v ^# | awk '{print $2}')";
echo "超时退出时间为:$(echo $TMOUT)"
echo "密码数字位数要求为:$(cat /etc/pam.d/sshd|grep difok | awk -F'dcredit=' '{print $2}' | awk '{print $1}')"
echo "密码小写位数要求为:$(cat /etc/pam.d/sshd|grep difok | awk -F'lcredit=' '{print $2}' | awk '{print $1}')"
echo "密码大写位数要求为:$(cat /etc/pam.d/sshd|grep difok | awk -F'ucredit=' '{print $2}' | awk '{print $1}')"
echo "密码特殊符号位数要求为:$(cat /etc/pam.d/sshd|grep difok | awk -F'ocredit=' '{print $2}' | awk '{print $1}')"
echo "密码长度要求为:$(cat /etc/pam.d/sshd|grep minlen | awk -F'dcredit=' '{print $2}' | awk '{print $1}')"
echo "登陆错误失败次数为:$(cat /etc/pam.d/sshd | grep deny | awk -F'deny=' '{print $2}' | awk '{print $1}')"
echo "失败锁定时常为:$(cat /etc/pam.d/sshd |grep unlock_time | awk -F'unlock_time=' '{print $2}' | awk '{print $1}')";
echo "是否限制Root用户远程登陆:$(cat /etc/ssh/sshd_config |grep -v "^#" |grep "PermitRootLogin" |awk -F' ' '{print $2}')（输出no为已限制）"
DENYROOT=$(cat /etc/pam.d/sshd |grep -v "^#"| grep even_deny_root |wc -l);
echo "是否设置ROOT用户失败锁定:$(if [ $DENYROOT -eq 1 ];then echo "yes";else echo "no";fi)"
echo "ROOT用户失败锁定时常为:$(cat /etc/pam.d/sshd |grep root_unlock_time | awk -F'root_unlock_time=' '{print $2}' | awk '{print $1}')"
echo "------------------------------------------------------------------------------------------------------umask信息"
echo "umask为:$(umask)"
echo "------------------------------------------------------------------------------------------------------系统安全策略信息"
echo "hosts.allow限制为:";if [ -z $(cat /etc/hosts.allow|grep -v "^#"|grep -v "^$") ];then echo "未做限制";else echo "$(cat /etc/hosts.allow|grep -v "^#"|grep -v "^$")";fi
echo "hosts.deny限制为:";if [ -z $(cat /etc/hosts.deny|grep -v "^#"|grep -v "^$") ];then echo "未做限制";else echo "$(cat /etc/hosts.deny|grep -v "^#"|grep -v "^$")";fi
echo "Selinux状态:$(cat /etc/selinux/config |grep -v "^#"|grep -v "^$"|awk -F 'SELINUX=' '{print $2}')"
Firewalld=$(systemctl status firewalld 1>/dev/null 2>/dev/null);
echo "Firewalld状态为:$(if [ $? -eq 0 ];then echo "yes";else echo "no";fi 2>/dev/null)"
iptabing=$(systemctl status iptables.service 1>/dev/null 2>/dev/null);
echo "iptables状态为:$(if [ $? -eq 0 ];then echo "yes";else echo "no";fi 2>/dev/null)"echo "防火墙策略为:$(iptables -L)"
echo "------------------------------------------------------------------------------------------------------开放端口信息"
echo "开放端口为:$(echo;netstat -anpt |grep LISTEN)"
echo "------------------------------------------------------------------------------------------------------rsys配置信息"
echo "rsys配置为:$(echo;cat /etc/rsyslog.conf  |grep -v -E "^#"|grep -v "^$")"
echo
auditip=$(cat /etc/rsyslog.conf  |grep -v -E "^#"|grep -v "^$"|grep "[0-9]\{1,3\}[.][0-9]\{1,3\}[.][0-9]\{1,3\}[.][0-9]\{1,3\}"|awk -F "@{1,2}" '{print $2}' |awk -F ":" '{print $1}');if [ "$auditip" != "" ] && [ "$auditip" != "127.0.0.1" ]; then ping $auditip -c 1 -W 1 >/dev/null;if [ $? -eq 0 ];then echo "审计日志服务器:" $auditip "通信正常"; else echo "审计日志服务器:" $auditip"通信失败"; fi; else echo "未配置审计服务地址";fi
echo "------------------------------------------------------------------------------------------------------重要文件权限信息"
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
echo "------------------------------------------------------------------------------------------------------定时任务"
echo "定时任务为:";cat /etc/crontab |grep -v "^#" |grep -v "^$"
echo "------------------------------------------------------------------------------------------------------系统日志"
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
echo "------------------------------------------------------------------------------------------------------"
echo "查看后三行secure日志:";head -n 3 /var/log/secure
echo "------------------------------------------------------------------------------------------------------登陆记录"
echo "查看用户登录记录:";last;
echo "------------------------------------------------------------------------------------------------------上次登录信息"
echo "查看所有用户上次登录信息：";lastlog`
	return cmd
}
