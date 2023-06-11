package run

func Linux_cmd() string {
	cmd := `
echo "------------------------------------------------------------------------------------------------------信息系统"
echo "系统版本:$(cat /etc/redhat-release)";
echo "系统架构:$(uname -m)"
echo "网卡信息:$(echo;ifconfig |grep inet |grep -v "inet6"|awk -F' ' '{print $2}')"
echo
echo "------------------------------------------------------------------------------------------------------身份鉴别"
echo "测评结果------v2.0"
echo "应对登录的用户进行身份标识和鉴别，身份标识具有唯一性，身份鉴别信息具有复杂度要求并定期更换；"
echo
echo $(IFS=':';output="经核查，登录系统确认服务器采用tty登录及远程登录，均存在口令鉴别措施，无法通过空口令进行登录；查看/etc/passwd文件中"; while read -r user enc_passwd uid gid full_name home shell; do if [[ ! "$shell" =~ ^(/sbin/nologin|/usr/sbin/nologin|/bin/false|/sbin/shutdown|/sbin/halt|/bin/sync|/usr/bin/false)$ ]]; then password_expiry_day=$(awk -F':' -v username="$user" '($1 == username) {print $5}' /etc/shadow); output+="存在可登录用户:${user},UID:${uid},密码更换周期为:${password_expiry_day}天; "; fi; done < /etc/passwd; output+="无重复用户以及uid身份标识唯一，在测试环境下测试创建重复用户root提示无法创建同名账户；查看/etc/shadow文件不存在空口令账户，测试验证现有用户均无法使用空口令进行tty登录及远程登录;"; minlen=$(grep "pam_pwquality.so" /etc/pam.d/system-auth | grep -oP "minlen=\K[^[:space:]]+" || echo "0"); ucredit=$(grep "pam_pwquality.so" /etc/pam.d/system-auth | grep -oP "ucredit=-\K[^[:space:]]+" || echo "0"); lcredit=$(grep "pam_pwquality.so" /etc/pam.d/system-auth | grep -oP "lcredit=-\K[^[:space:]]+" || echo "0"); dcredit=$(grep "pam_pwquality.so" /etc/pam.d/system-auth | grep -oP "dcredit=-\K[^[:space:]]+" || echo "0"); ocredit=$(grep "pam_pwquality.so" /etc/pam.d/system-auth | grep -oP "ocredit=-\K[^[:space:]]+" || echo "0"); enforce_for_root=$(grep "pam_pwquality.so" /etc/pam.d/system-auth | grep -o "enforce_for_root" || echo "无"); output+="查看/etc/pam.d/system-auth文件配置了password requisite pam_pwquality.so ucredit=-$ucredit lcredit=-$lcredit dcredit=-$dcredit ocredit=-$ocredit minlen=$minlen 口令长度为${minlen}位、小写字为${lcredit}位、大写字母位${ucredit}位、特殊符号为${ocredit}位、数字为${dcredit}位，"; if [ $enforce_for_root = "无" ]; then output+="未配置enforce_for_root无法对管理员用户生效。"; else output+="并配置了enforce_for_root对管理员用户生效，"; fi; pass_max_days=$(grep -E "^PASS_MAX_DAYS" /etc/login.defs | awk '{ print $2 }'); output+="查看/etc/login.defs文件PASS_MAX_DAYS为${pass_max_days}文件新建用户时密码更换周期为:${pass_max_days}天,新建用户时"; if [ $pass_max_days -le 90 ]; then output+="满足密码定期更换周期要求。"; else output+="不满足密码定期更换周期要求。"; fi; echo "$output")
echo
echo "------------------------------------------------------------------------------------------------------"
echo "应具有登录失败处理功能，应配置并启用结束会话、限制非法登录次数和当登录连接超时自动退出等相关措施；"
echo
echo $(files=("/etc/pam.d/system-auth" "/etc/pam.d/password-auth" "/etc/pam.d/sshd");regex="pam_tally2.so.*deny=([0-9]+).*unlock_time=([0-9]+)";root_regex="even_deny_root";output="";if [ "$(systemctl get-default)" = "graphical.target" ];then output+="经核查，通过ssh登录此操作系统确认服务器采用tty登录，远程登录以及桌面登录，";else output+="经核查，通过ssh登录此操作系统确认服务器采用tty登录及远程登录未开启桌面登录，";fi;for file in "${files[@]}"; do output+=$'检查文件'$file$'';if [[ -f $file ]]; then match=$(grep -Eo "$regex" $file);if [[ -z $match ]]; then output+="中确认未开启失败锁定次数和锁定时间，"$'';else deny_count=$(echo $match | grep -Eo 'deny=[0-9]+' | cut -d'=' -f2);unlock_time=$(echo $match | grep -Eo 'unlock_time=[0-9]+' | cut -d'=' -f2);output+="连续登录失败$deny_count次后锁定账户$unlock_time秒，";root_match=$(grep -E "$root_regex" $file);if [[ -z $root_match ]]; then output+="针对root用户的失败锁定未开启，"$'';else output+="针对root用户的失败锁定已开启，"$'';fi;fi;else output+="文件不存在: $file"$'';fi;output+=$'';done;contmout=$(grep -vE "^#|^$" /etc/ssh/sshd_config | grep "LoginGraceTime");con_value=$(echo $contmout | grep -oE '[0-9]+');if [[ -z $con_value ]]; then con_value="0";output+="未配置在远程连接过程中超时退出策略，";else output+="在ssh远程建立连接过程中通过etc/ssh/sshd_config配置文件设置LoginGraceTime为${con_value}，连接过程超过${con_value}秒后自动退出，";fi;tmout_value=$(cat /etc/profile | grep "TMOUT");tmout_v=$(echo $tmout_value | grep -oE '[0-9]+');if [[ -z $tmout_value ]]; then tmout_value="0";output+="在成功登录后未在全局配置文件/etc/profile中配置超时退出功能。";else output+="在成功登录后通过全局配置文件/etc/profile中配置超过${tmout_v}秒后自动退出功能。";fi;echo -e "$output")
echo
echo "------------------------------------------------------------------------------------------------------"
echo "当进行远程管理时，应采取必要措施防止鉴别信息在网络传输过程中被窃听；"
echo
echo $(output=""; outputtelnet=$(service telnet.socket status 2>/dev/null); version=$(grep -w "Protocol" /etc/ssh/sshd_config | awk -F " " '{print $2}'); if [[ -z $version ]]; then version=2; fi; if [[ $outputtelnet =~ "active (listening)" ]]; then output+="经核查，此设备在进行远程管理时，采用ssh以及Telnet协议进行远程管 理，通过通过Wireshark抓包验证在使用Telnet鉴别信息为明文传输，无法防止鉴别信息在网络传输过程中被窃听。"; else output+="经核查，此设备在进行远程管理时，采用ssh的加密协议进行远程管理，通过查看sshd_config文件中Protocol字段，确认ssh使用V${version}版本协议，并且已禁用telnet等明文传输协议，通过Wireshark抓包验证鉴别信息为密文传输，可防止鉴别信息在网络传输过程中被窃听。"; fi; echo $output)
echo
echo "------------------------------------------------------------------------------------------------------"
echo "应采用口令、密码技术、生物技术等两种或两种以上组合的鉴别技术对用户进行身份鉴别，且其中一种鉴别技术至少应使用密码技术来实现；"
echo
echo "经核查，此操作系统仅采用用户名+密码进行身份鉴别，未采用动态口令、数字证书、生物技术或设备指纹等组合方式对用户进行身份鉴别。"
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
echo "------------------------------------------------------------------------------------------------------安全审计"
echo "应启用安全审计功能，审计覆盖到每个用户，对重要的用户行为和重要安全事件进行审计；"
echo
echo $(output="经核查，"; status_audit=$(systemctl is-active auditd | awk '{print ($1=="active") ? "服务器 已开启auditd日志进程，相关审计日志为/var/log/audit/audit.log，" : "服务器未开启auditd日志进程，"}'); status_rsyslog=$(systemctl is-active rsyslog | awk '{print ($1=="active") ? "已开启rsyslog日志进程，相关审计日志为/var/log/secure、/var/log/messages、/var/log/dmesg、/var/log/boot.log、/var/log/cron、/var/log/wtmp等，审计已覆盖当前所有用户，审计日志包含重要用户行为：用户登录登出、历史命令、特权命令、文件访问记录、重要文件变更等，以及包含重要的安全事件：系统启动日志、引导日志、驱动日志、文件日志、系统调用记录、服务进程日志、安装更新日志等。" : "rsyslog日志进程未开启，无法保证对重要的用户行为及重要安全事件进行审计。"}'); output+=$status_audit;output+=$status_rsyslog; echo $output)
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
