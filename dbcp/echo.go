package dbcp

import (
	"fmt"
)

func Huawei() {
	fmt.Println(`
	display version（查看版本）
	disp curr   （查看当前配置）
	disp vlan    （查看VLAN信息）
	disp interface brief （查看接口信息，带宽利用率）
	disp cpu-usage （查看cpu使用率）
	disp memory-usage （查看内存使用情况）有问题。
	disp acl all  （查看访问控制列表）
	disp password-control （查看密码策略，超时等）（华三）
	disp info-center  （查看日志审计开启情况、日志服务器配置）
	disp logbuffer  （查看日志审计的内容）
	disp interface	(接口信息)`)
}

func Dm() {
	fmt.Println(`
	默认端口:5236
	进入命令行界面：找到安装位置中的tool目录，执行./disql
	SQL> conn sysdba/gys123...@127.0.0.1:5236

	服务器[127.0.0.1:5236]:处于普通打开状态
	登录使用时间 : 2.131(ms)
	默认ui管理工具：
	找到安装位置中的tool目录，执行./manager

	审计日志工具，必须使用SYSAUDITOR用户
	找到安装位置中的tool目录，执行./analyzer，双击审计日志查看器，选择审计文件
	sql日志以“dmsql-实例名-时间-标号” 命名, 默认生成在 DM安装目录log 子目录下面
	SP_SET_ENABLE_AUDIT (1); 
	ENABLE_AUDIT=0 —关闭审计
	ENABLE_AUDIT=1 —打开普通审计
	ENABLE_AUDIT=2 —打开普通审计和实时审计
	查看审计状态
	SELECT * FROM V$DM_INI WHERE PARA_NAME='ENABLE_AUDIT';




	查看用户、有限期、策略
	SELECT * FROM V$PARAMETER WHERE NAME='PWD_POLICY';
	SELECT USERNAME,PASSWORD_VERSIONS,EXPIRY_DATE FROM DBA_USERS;
	0 无策略
	1 禁止与用户名相同
	2 口令长度不小于9
	4 至少包含一个大写字母（A-Z）
	8 至少包含一个数字（0-9）
	16 至少包含一个标点符号（英文输入法状态下，除“和空格外的所有符号）

	用户信息:
	SELECT * FROM sysusers;

	查看审计状态：
	SELECT * FROM V$DM_INI WHERE PARA_NAME='ENABLE_AUDIT';

	查看日志文件路径信息
	select path,rlog_size from v$rlogfile;

	查看是否开启ssl：
	SELECT NAME,TYPE,VALUE FROM V$PARAMETER WHERE NAME='ENABLE_ENCRYPT' OR NAME
	='COMM_ENCRYPT_NAME';
	ENABLE_ENCRYPT为0，代表不加密，配置为1代表ssl加密，2代表ssl认证；COMM_ENCRYPT_NAME值为空，未配置加密算法，

	查看资源信息基本测评项都包括了：
	SELECT B.NAME,    -- 用户名称
	-- A.ID,				-- 用户 ID
	-- A.PASSWORD,			-- 用户口令
	A.AUTHENT_TYPE,		-- 用户认证方式：NDCT_DB_AUTHENT/NDCT_OS_AUTHENT/NDCT_NET_AUTHENT/NDCT_UNKOWN_AUTHENT
	A.SESS_PER_USER,  -- 在一个实例中，一个用户可以同时拥有的会话数量
	A.CONN_IDLE_TIME,	-- 用户会话的最大空闲时间
	A.FAILED_NUM,	-- 用户登录失败次数限制
	A.LIFE_TIME, 	-- 一个口令在终止使用前可以使用的天数
	A.REUSE_TIME,	-- 一个口令在可以重新使用之前必须经过的天数
	A.REUSE_MAX,	-- 一个口令在可以重新使用前必须改变的次数
	A.LOCK_TIME,	-- 用户口令锁定时间
	A.GRACE_TIME,	-- 用户口令过期后的宽限时间
	A.LOCKED_STATUS,	-- 用户登录是否锁定：LOGIN_STATE_UNLOCKED/LOGIN_STATE_LOCKED
	A.LASTEST_LOCKED,	-- 用户最后一次的锁定时间
	A.PWD_POLICY,	-- 用户口令策略：NDCT_PWD_POLICY_NULL/NDCT_PWD_POLICY_1/NDCT_PWD_POLICY_2/NDCT_PWD_POLICY_3/NDCT_PWD_POLICY_4/NDCT_PWD_POLICY_5
	A.RN_FLAG,	-- 是否只读
	A.ALLOW_ADDR,	-- 允许的 IP 地址
	A.NOT_ALLOW_ADDR,	-- 不允许的 IP 地址
	A.ALLOW_DT,		-- 允许登录的时间段
	A.NOT_ALLOW_DT,		-- 不允许登录的时间段
	A.LAST_LOGIN_DTID,	-- 上次登录时间
	A.LAST_LOGIN_IP,	-- 上次登录 IP 地址
	A.FAILED_ATTEMPS,	-- 将引起一个账户被锁定的连续注册失败的次数
	A.ENCRYPT_KEY		-- 用户登录的存储加密密钥
	FROM SYSUSERS A,SYS.SYSOBJECTS B 
	WHERE A.ID=B.ID;


	`)

}

func Cisco() {
	fmt.Println(`
		查看系统版本及设备型号
		show version	
		查看当前配置
		show running-config	
		查看VLAN信息
		show vlan  brief
		查看接口信息，带宽利用率
		show interfaces status		
		查看cpu使用率
		show processes
		查看cpu使用率
		show processes cpu	
		查看内存使用率
		show processes memory	
		查看访问控制列表
		show  access-lists
		查看日志审计开启情况、日志服务器配置
		show logging
		查看日志审计的内容
		show logbuffer`)
}

func Linux() {
	fmt.Println(`
	lsb_release -a  
	cat /etc/redhat-release 
	系统版本

	cat /etc/passwd
	cat /etc/shadow
	账户名称、特权用户

	cat /etc/login.defs |grep -v '#' |grep -v "^$"
	密码有效期

	grep -v '#' /etc/pam.d/login | grep -v '^$'
	登录失败重试次数、锁定策略

	grep -v '#' /etc/pam.d/system-auth | grep -v '^$'
	密码复杂度策略

	grep -v '#' /etc/ssh/sshd_config  | grep -v '^$'
	grep -v '#' /etc/pam.d/sshd  | grep -v '^$'
	登录失败重试次数、锁定策略

	sestatus -v
	selinux状态、对重要主体和客体设置安全标记，通过安全标记控制主体对客体的访问

	set |gerp "TMOUT"
	cat /etc/profile|greo "TMOUT"
	超时策略

	ps -ef |  grep rsyslog
	grep -v '#' /etc/rsyslog.conf   | grep -v '^$'
	日志审计 rsyslog状态

	ps -ef |  grep audit
	日志审计 audit状态

	auditctl -s 
	auditctl -l
	审计信息
	
	grep -v '#' /etc/hosts.allow   | grep -v '^$'
	grep -v '#' /etc/hosts.deny   | grep -v '^$'
	地址限制

	iptables -L
	防火墙状态
	ls -l /var/log/audit/audit.log /var/log/secure /var/log/messages | cut -d ' ' -f1
	netstat -tuanl
	开放端口
	crontab -l
	计划任务、备份策略
	`)

}

func Mysql() {
	fmt.Println(`版本信息:
		select version()
	
		用户名、主机、密码有效期，远程登录:
		select user,host,password_expired from mysql.user;
	
		查看密码策略：
		show variables like 'validate_password%';
	
		查看最大连接数和状态：
		show variables like '%connect%';
	
		查看超时时间
		show variables like '%timeout%';
	
		查看登录失败的次数及延迟响应时间
		show variables like '%connection_control%';
	
		查看安装的插件:
		show plugins;
	
		查看是否开启ssl:
		show global variables like '%ssl%';
		
		确定用户是否可以关闭MySQL服务器
		确定用户是否可以执行SELECT INTO OUTFILE和LOAD DATA INFILE命令
		确定用户是否可以将已经授予给该用户自己的权限再授予其他用户
		select Host,User,File_priv,Shutdown_priv,grant_priv from mysql.user;
	
		查看是否日志审计
		show variables like 'log_%';
	
		查看是否开启审计:
		show global variables like '%general%';
		
		查看节点信息:
		show master status `)
}

func Aix() {
	fmt.Println(` #版本
	oslevel   
	oslevel -r
	instfix -i | grep AIX_ML
	grep -v '#' /etc/security/passwd  | grep -v '^$'| grep -v '*'
	pwck
	grep -v '#' /etc/hosts.equiv  | grep -v '^$'| grep -v '*'
	grep -v '#' /etc/security/user  | grep -v '^$'| grep -v '*'
	grep -v '#' /etc/security/login.cfg  | grep -v '^$' | grep -v '*'
	grep -v '#' /etc/profile  | grep -v '^$' | grep -v '*'
	grep -v '#' /etc/ssh/ssh_config  | grep -v '^$' | grep -v '*'
	ps -elf | grep ftp
	ps -elf | grep ssh
	ps -elf | grep telnet
	ls -l /etc/passwd   /etc/security/passwd /var/log/messages
	lssrc -s syslogd
	ps -ef | grep syslog
	grep -v '#' /etc/syslog.conf | grep -v '^$' | grep -v '*'
	audit query
	grep -v '#' /etc/security/audit/events  | grep -v '^$' | grep -v '*'
	grep -v '#' /etc/security/audit/config | grep -v '^$' | grep -v '*'
	grep -v '#' /etc/security/limits  | grep -v '^$' | grep -v '*'
	开放端口
	netstat -tuan
	`)

}

func Postsql() {
	fmt.Println(`
	命令行进入psql -U username -h ipaddress -d dbname
	默认端口:5432

	查看版本
	select version();
	
	查看密码复杂度：
	show shared_preload_libraries;
	
	密码定期更换：查看valuntil字段
	select * from pg_shadow;
	
	查看ssl:
	cat $PGDATA/postgresql.conf |grep ssl
	
	查看存储加密算法:
	show password_encryption;
	
	查看审计状态：
	show logging_collector;
	
	查看审计类型：  --日志记录类型，默认是stderr，只记录错误输出，推荐csvlog，总共包含：stderr, csvlog, syslog, and eventlog。
	show log_destination;`)

}

func Oracle() {
	fmt.Println(`
	登录用户：
	show user;
	select username,account_status  from dba_users where profile = 'DEFAULT';

	失败次数
	select * from dba_profiles s where s.profile='DEFAULT' and resource_name='FAILED_LOGIN_ATTEMPTS';

	密码定期更换时间：
	select * from dba_profiles s where s.profile='DEFAULT' AND resource_name='PASSWORD_LIFE_TIME';
	
	密码复杂度
	select * from dba_profiles s where s.profile='DEFAULT' and resource_name='Password_verify_function';
	
	超时时间
	select * from dba_profiles s where s.profile='DEFAULT' and resource_name='IDLE_TIME';
	
	查看哪些用户有sysdba或sysoper系统权限(查询时需要相应权限)
	select * from V$PWFILE_USERS;
	
	查看是否开启审计#VALUE值为DB，表面审计功能为开启的状态
	show parameter audit_trail;
	
	查看审计数据：
	select * from aud$;
	select * FROM SYS.AUD$;
	
	查看安全标记：
	SELECT VALUE FROM V$OPTION WHERE PARAMETER = 'Oracle Label Security';
	`)
}
