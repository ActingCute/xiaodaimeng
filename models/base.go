package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"xiaodaimeng/public"
)

var Mysql *xorm.Engine
var DBOk = false

var SQL = public.ConfigData.Mysql

func InitMysql() {

	var err error
	dataSourceName := SQL.User + ":" + SQL.Password + "@tcp(" + SQL.Host + ":" + SQL.Port + ")/" + SQL.Table + "?charset=utf8&parseTime=true&loc=Local"

	Mysql, err = xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		public.Error("数据库连接失败:", err)
		return
	} else {
		public.Debug("数据库连接ok:")
	}
}
