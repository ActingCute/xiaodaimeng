package models

import (
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"

	"xiaodaimeng/public"
)

var DBOk = false
var DB = new(xorm.Engine)

func InitDB() {
	public.Debug("连接数据表")
	var err error
	DB, err = xorm.NewEngine("sqlite3", "./xiaodaimeng.db")
	if err != nil {
		public.Error(err)
		return
	}

	err = DB.Sync2(new(Work))

	if err != nil {
		public.Error(err)
		return
	}

	work := Work{
		WxId:  "test",
		Type:  "type",
		Msg:   "okk",
		Other: "pppp",
	}
	InsertWork(&work)
	public.Debug(work)
}
