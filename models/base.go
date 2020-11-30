package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"os"

	"xiaodaimeng/public"
)

var DBOk = false
var DB = new(xorm.Engine)
var SQL = public.ConfigData.Mysql
var ShowSql = public.ConfigData.IsDebug

const (
	TableBlackList = "black_list"
	TableWork = "work"
)

func initMysql() {

	var err error
	dataSourceName := SQL.User + ":" + SQL.Password + "@tcp(" + SQL.Host + ":" + SQL.Port + ")/" + SQL.Table + "?charset=utf8&parseTime=true&loc=Local"

	DB, err = xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		public.Error("数据库连接失败:", err)
		return
	}
	public.Debug("mysql数据库连接ok")

	err = DB.Sync2(new(Work),new(BlackList))

	if err != nil {
		public.Error("InitDB:",err)
		return
	}
	DBOk = true

	DB.ShowSQL(ShowSql)
}

func initSqlite3()  {
	public.Debug("连接数据表")
	var err error
	DB, err = xorm.NewEngine("sqlite3", "./xiaodaimeng.db")
	if err != nil {
		public.Error(err)
		return
	}

	err = DB.Sync2(new(Work),new(BlackList))

	if err != nil {
		public.Error("InitDB:",err)
		return
	}
	public.Debug("sqlite3数据库连接ok")

	DBOk = true
}

func InitDB() {
	if public.ConfigData.UseMysql {
		//使用mysql
		initMysql()
	}else{
		//使用sqlite3
		initSqlite3()
	}
	if DBOk {
		InitBlacklist()
		f, err1 := os.Create("sql.log")
		if err1 != nil {
			public.Debug(err1.Error())
		}
		DB.SetLogger(xorm.NewSimpleLogger(f))
		defer f.Close()
	}
}
