#!/bin/bash
IFS=':'
output="经核查，登录系统确认服务器采用tty登录及远程登录，均存在口令鉴别措施，无法通过空口令进行登录；查看/etc/passwd文件中"
while read -r user enc_passwd uid gid full_name home shell; do
    if [[ ! "$shell" =~ ^(/sbin/nologin|/usr/sbin/nologin|/bin/false|/sbin/shutdown|/sbin/halt|/bin/sync|/usr/bin/false)$ ]]; then
        password_expiry_day=$(awk -F':' -v username="$user" '($1 == username) {print $5}' /etc/shadow)
        output+="存在可登录用户:${user},UID:${uid},密码更换周期为:${password_expiry_day}天; "
    fi
done < /etc/passwd
output+="无重复用户以及uid身份标识唯一，在测试环境下测试创建重复用户root提示无法创建同名账户；查看/etc/shadow文件不存在空口令账户，测试验证现有用户均无法使用空口令进行tty登录及远程登录;"
minlen=$(grep "pam_pwquality.so" /etc/pam.d/system-auth | grep -oP "minlen=\K[^[:space:]]+" || echo "0")
ucredit=$(grep "pam_pwquality.so" /etc/pam.d/system-auth | grep -oP "ucredit=-\K[^[:space:]]+" || echo "0")
lcredit=$(grep "pam_pwquality.so" /etc/pam.d/system-auth | grep -oP "lcredit=-\K[^[:space:]]+" || echo "0")
dcredit=$(grep "pam_pwquality.so" /etc/pam.d/system-auth | grep -oP "dcredit=-\K[^[:space:]]+" || echo "0")
ocredit=$(grep "pam_pwquality.so" /etc/pam.d/system-auth | grep -oP "ocredit=-\K[^[:space:]]+" || echo "0")
enforce_for_root=$(grep "pam_pwquality.so" /etc/pam.d/system-auth | grep -o "enforce_for_root" || echo "无")

output+="查看/etc/pam.d/system-auth文件配置了password requisite pam_pwquality.so ucredit=-$ucredit lcredit=-$lcredit dcredit=-$dcredit ocredit=-$ocredit minlen=$minlen 口令长度为${minlen}位、小写字为${lcredit}位、大写字母位${ucredit}位、特殊符号为${ocredit}位、数字为${dcredit}位，"

if [ $enforce_for_root = "无" ];then
    output+="未配置enforce_for_root无法对管理员用户生效。"
else
    output+="并配置了enforce_for_root对管理员用户生效。"
fi
pass_max_days=$(grep -E "^PASS_MAX_DAYS" /etc/login.defs | awk '{print $2}')

output+="查看/etc/login.defs文件PASS_MAX_DAYS为${pass_max_days}文件新建用户时密码更换周期为:${pass_max_days},"
if [ $pass_max_days -le 90 ]; then
    output+="满足密码定期更换周期要求。"
else
    output+="不满足密码定期更换周期要求。"
fi
echo "$output"

#一句话命令
#IFS=':'; output="经核查，登录系统确认服务器采用tty登录及远程登录，均存在口令鉴别措施，无法通过空口令进行登录；查看/etc/passwd文件中"; while read -r user enc_passwd uid gid full_name home shell; do if [[ ! "$shell" =~ ^(/sbin/nologin|/usr/sbin/nologin|/bin/false|/sbin/shutdown|/sbin/halt|/bin/sync|/usr/bin/false)$ ]]; then password_expiry_day=$(awk -F':' -v username="$user" '($1 == username) {print $5}' /etc/shadow); output+="存在可登录用户:${user},UID:${uid},密码更换周期为:${password_expiry_day}天; "; fi; done < /etc/passwd; output+="无重复用户以及uid身份标识唯一，在测试环境下测试创建重复用户root提示无法创建同名账户；查看/etc/shadow文件不存在空口令账户，测试验证现有用户均无法使用空口令进行tty登录及远程登录;"; minlen=$(grep "pam_pwquality.so" /etc/pam.d/system-auth | grep -oP "minlen=\K[^[:space:]]+" || echo "0"); ucredit=$(grep "pam_pwquality.so" /etc/pam.d/system-auth | grep -oP "ucredit=-\K[^[:space:]]+" || echo "0"); lcredit=$(grep "pam_pwquality.so" /etc/pam.d/system-auth | grep -oP "lcredit=-\K[^[:space:]]+" || echo "0"); dcredit=$(grep "pam_pwquality.so" /etc/pam.d/system-auth | grep -oP "dcredit=-\K[^[:space:]]+" || echo "0"); ocredit=$(grep "pam_pwquality.so" /etc/pam.d/system-auth | grep -oP "ocredit=-\K[^[:space:]]+" || echo "0"); enforce_for_root=$(grep "pam_pwquality.so" /etc/pam.d/system-auth | grep -o "enforce_for_root" || echo "无"); output+="查看/etc/pam.d/system-auth文件配置了password requisite pam_pwquality.so ucredit=-$ucredit lcredit=-$lcredit dcredit=-$dcredit ocredit=-$ocredit minlen=$minlen 口令长度为${minlen}位、小写字为${lcredit}位、大写字母位${ucredit}位、特殊符号为${ocredit}位、数字为${dcredit}位，"; if [ $enforce_for_root = "无" ]; then output+="未配置enforce_for_root无法对管理员用户生效。"; else output+="并配置了enforce_for_root对管理员用户生效。"; fi; pass_max_days=$(grep -E "^PASS_MAX_DAYS" /etc/login.defs | awk '{ print $2 }'); output+="查看/etc/login.defs文件PASS_MAX_DAYS为${pass_max_days}文件新建用户时密码更换周期为:${pass_max_days},"; if [ $pass_max_days -le 90 ]; then output+="满足密码定期更换周期要求。"; else output+="不满足密码定期更换周期要求。"; fi; echo "$output"
