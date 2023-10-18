package run

import (
	"database/sql"
	"embed"
)

//go:embed oracle_html.html
var templateFileOracle embed.FS

type DataOracle struct {
	Name          string
	Version       []string    //版本信息
	UserInfo      []SysUser   //sys.user视图数据
	ListParameter []Parameter //部分安全配置策略
	DBAUSERS      []DBA_USERS //DBA_USERS信息
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
	User    string       `db:"USERNAME"`
	Profile string       `db:"PROFILE"`
	Status  string       `db:"ACCOUNT_STATUS"`
	Expiry  sql.NullTime `db:"EXPIRY_DATE"`
}

type Parameter struct {
	Name  string `db:"NAME"`
	Value string `db:"VALUE"`
}
