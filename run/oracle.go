package run

import (
	"database/sql"
	"fmt"
	_ "github.com/sijms/go-ora/v2"
	"github.com/spf13/cobra"
	"golin/global"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

func Oraclestart(cmd *cobra.Command, args []string) {
	//确认结果是否输出
	echotype, err := cmd.Flags().GetBool("echo")
	if err != nil {
		fmt.Println(err)
		return
	}
	echorun = echotype

	//获取分隔符
	spr, err := cmd.Flags().GetString("spript")
	if err != nil {
		fmt.Println(err)
		return
	}
	//如果value值不为空则是运行一次的模式
	value, err := cmd.Flags().GetString("value")
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(value) > 10 {
		Onlyonerun(value, spr, "oracle")
		wg.Wait()
		return
	}
	// 下面开始执行批量的
	ippath, err := cmd.Flags().GetString("ip")
	if err != nil {
		fmt.Println(err)
		return
	}
	//判断pgsql.txt文件是否存在
	Checkfile(ippath, fmt.Sprintf("名称%sip%s用户%s密码%s端口", Split, Split, Split, Split), global.FilePer, ippath)
	// 运行share文件中的函数
	Rangefile(ippath, spr, "oracle")
	wg.Wait()
	//完成前最后写入文件
	Deffile("oracle", count, count-len(errhost), errhost)
}

func OracleRun(name, host, user, passwd, port string) error {
	defer wg.Done()
	sid := "orcl"
	oidlist := strings.Split(name, "sid=")
	if len(oidlist) == 2 {
		sid = oidlist[1]
	}
	dsn := fmt.Sprintf("oracle://%s:%s@%s:%s/%s?timeout=15", user, passwd, host, port, sid)
	db, err := sql.Open("oracle", dsn)
	if err != nil {
		return fmt.Errorf("%s:连接失败！检查网络或用户密码或SID吧！", host)
	}

	fullPath := filepath.Join(succpath, "oracle")
	firenmame := filepath.Join(fullPath, fmt.Sprintf("%s_%s.html", name, host))
	//先删除之前的同名记录文件
	os.Remove(firenmame)

	data := DataOracle{}
	data.Name = fmt.Sprintf("%s_%s", name, host)
	//版本信息
	banners, err := QueryAndParse(db, "select BANNER from v$version", func(rows *sql.Rows) (string, error) {
		var banner string
		err := rows.Scan(&banner)
		return banner, err
	})
	if err != nil && len(banners) == 0 {
		return fmt.Errorf("%s:执行version命令失败,检查网络或用户密码或SID吧！%v", host, err)
	}
	data.CSS = css
	data.Version = banners

	//读取sys.user的视图信息
	users, _ := QueryAndParse(db, "SELECT USER#, NAME, TYPE#, PASSWORD, CTIME, PTIME, EXPTIME, LTIME FROM sys.user$", func(rows *sql.Rows) (SysUser, error) {
		var userinfo SysUser
		err := rows.Scan(&userinfo.UserNum, &userinfo.Name, &userinfo.TypeNum, &userinfo.Password, &userinfo.CTime, &userinfo.PTime, &userinfo.ExpTime, &userinfo.LTime)
		return userinfo, err
	})
	data.UserInfo = users

	//v$parameter部分安全配置
	par, _ := QueryAndParse(db, `SELECT NAME, VALUE FROM v$parameter WHERE NAME IN (
        'remote_login_passwordfile', 
        'remote_os_authent', 
        'sec_case_sensitive_logon', 
        'audit_trail'
    )`, func(rows *sql.Rows) (Parameter, error) {
		var info Parameter
		err := rows.Scan(&info.Name, &info.Value)
		return info, err
	})
	data.ListParameter = par

	//DBA_USERS部分信息
	dbauser, _ := QueryAndParse(db, `SELECT USERNAME,PROFILE,ACCOUNT_STATUS,EXPIRY_DATE,LOCK_DATE,CREATED FROM DBA_USERS`, func(rows *sql.Rows) (DBA_USERS, error) {
		var info DBA_USERS
		err := rows.Scan(&info.User, &info.Profile, &info.Status, &info.Expiry, &info.LockTime, &info.CreateTime)
		return info, err
	})
	data.DBAUSERS = dbauser

	//确认配置文件的安全策略规则
	Passwdverify, _ := QueryAndParse(db, `SELECT * FROM DBA_PROFILES`, func(rows *sql.Rows) (VerifyFunc, error) {
		var info VerifyFunc
		err := rows.Scan(&info.Profile, &info.Resp, &info.Type, &info.Limit)
		return info, err
	})
	data.PasswdVerify = Passwdverify

	//用户系统权限
	systemuser, _ := QueryAndParse(db, `SELECT u.USERNAME, sp.PRIVILEGE, sp.ADMIN_OPTION
FROM DBA_USERS u
JOIN DBA_SYS_PRIVS sp ON u.USERNAME = sp.GRANTEE
WHERE u.ACCOUNT_STATUS = 'OPEN'`, func(rows *sql.Rows) (AUTHORITY, error) {
		var info AUTHORITY
		err := rows.Scan(&info.Name, &info.Privilege, &info.Opthion)
		return info, err
	})
	data.SYSTEMAUTHORITY = systemuser

	//用户对象权限
	Objectuser, _ := QueryAndParse(db, `SELECT u.USERNAME, tp.TABLE_NAME, tp.PRIVILEGE
FROM DBA_USERS u
JOIN DBA_TAB_PRIVS tp ON u.USERNAME = tp.GRANTEE
WHERE u.ACCOUNT_STATUS = 'OPEN'`, func(rows *sql.Rows) (AUTHORITY, error) {
		var info AUTHORITY
		err := rows.Scan(&info.Name, &info.Privilege, &info.Opthion)
		return info, err
	})
	data.ObjectPermissions = Objectuser

	//审计日志配置
	audit, _ := QueryAndParse(db, `SELECT NAME,VALUE,ISSES_MODIFIABLE,ISSYS_MODIFIABLE,ISINSTANCE_MODIFIABLE,DESCRIPTION
FROM V$PARAMETER
WHERE NAME LIKE '%audit%'`, func(rows *sql.Rows) (Audit, error) {
		var info Audit
		err := rows.Scan(&info.NAME, &info.VALUE, &info.ISSES_MODIFIABLE, &info.ISSYS_MODIFIABLE, &info.ISINSTANCE_MODIFIABLE, &info.DESCRIPTION)
		return info, err
	})
	data.AuditPARAMETER = audit

	//超时
	timeout, _ := QueryAndParse(db, `SELECT PROFILE, RESOURCE_NAME, LIMIT
FROM dba_profiles
WHERE resource_name = 'IDLE_TIME'`, func(rows *sql.Rows) (TimeoutidleTime, error) {
		var info TimeoutidleTime
		err := rows.Scan(&info.Profile, &info.ResourceNam, &info.Limit)
		return info, err
	})
	data.IdleTime = timeout

	//所有函数
	funcpass, _ := QueryAndParse(db, `SELECT OWNER,NAME,LINE,TEXT FROM dba_source WHERE TYPE= 'FUNCTION'`, func(rows *sql.Rows) (DbaSource, error) {
		var info DbaSource
		err := rows.Scan(&info.OWNER, &info.NAME, &info.LINE, &info.TEXT)
		return info, err
	})
	data.FuncPass = funcpass

	// 读取模板文件
	tmpl, err := template.ParseFS(templateFileOracle, "template/oracle_html.html")
	if err != nil {
		return fmt.Errorf("%s:读取HTML模板文件失败。", host)
	}

	//确认采集完成目录是否存在
	_, err = os.Stat(fullPath)
	if os.IsNotExist(err) {
		os.MkdirAll(fullPath, os.FileMode(global.FilePer))
	}

	// 创建一个新的文件
	newFile, err := os.Create(firenmame)
	if err != nil {
		return fmt.Errorf("%s:创建文件失败。", host)
	}
	defer newFile.Close()
	// 将模板执行的结果写入新的文件
	err = tmpl.Execute(newFile, data)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("%s:保存模板文件失败。", host)
	}
	return nil

}

func QueryAndParse[T any](db *sql.DB, query string, parseRow func(*sql.Rows) (T, error)) ([]T, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		result, err := parseRow(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}
