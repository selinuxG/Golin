package run

import (
	"database/sql"
	"html/template"
)

type DataOracle struct {
	CSS               template.CSS      //css样式
	Name              string            //名称
	Version           []string          //版本信息
	UserInfo          []SysUser         //sys.user视图数据
	ListParameter     []Parameter       //部分安全配置策略
	DBAUSERS          []DBA_USERS       //DBA_USERS信息
	PasswdVerify      []VerifyFunc      //配置文件的密码策略规则
	SYSTEMAUTHORITY   []AUTHORITY       //系统权限
	ObjectPermissions []AUTHORITY       //对象权限
	AuditPARAMETER    []Audit           //审计配置
	IdleTime          []TimeoutidleTime //超时
	FuncPass          []DbaSource       //数据库中所有函数
}

type SysUser struct {
	UserNum  sql.NullInt64  `db:"USER#"`
	Name     sql.NullString `db:"NAME"`
	TypeNum  sql.NullInt64  `db:"TYPE#"`
	Password sql.NullString `db:"PASSWORD"`
	CTime    sql.NullTime   `db:"CTIME"`
	PTime    sql.NullTime   `db:"PTIME"`
	ExpTime  sql.NullTime   `db:"EXPTIME"`
	LTime    sql.NullTime   `db:"LTIME"`
}

type DBA_USERS struct {
	User       string       `db:"USERNAME"`
	Profile    string       `db:"PROFILE"`
	Status     string       `db:"ACCOUNT_STATUS"`
	Expiry     sql.NullTime `db:"EXPIRY_DATE"`
	LockTime   sql.NullTime `db:"LOCK_DATE"`
	CreateTime sql.NullTime `db:"CREATED"`
}

type Parameter struct {
	Name  string `db:"NAME"`
	Value string `db:"VALUE"`
}

type VerifyFunc struct {
	Profile string
	Resp    string
	Type    string
	Limit   string
}

type AUTHORITY struct {
	Name      string
	Privilege string
	Opthion   string
}

type Audit struct {
	NAME                  string
	VALUE                 sql.NullString `db:"VALUE"`
	ISSES_MODIFIABLE      string
	ISSYS_MODIFIABLE      string
	ISINSTANCE_MODIFIABLE string
	DESCRIPTION           string
}

type TimeoutidleTime struct {
	Profile     string
	ResourceNam string
	Limit       string
}

type DbaSource struct {
	OWNER string
	NAME  string
	LINE  int
	TEXT  string
}
