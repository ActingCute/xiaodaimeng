package models

import (
	//_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"xiaodaimeng/public"
)

var DBOk = false
var DB = new(sql.DB)

func InitDB() {
	public.Debug("打开数据表")
	var err error
	DB, err = sql.Open("sqlite3", "xiaodaimeng.db")
	if err != nil {
		public.Error(err)
		return
	}

	fmt.Println("生成数据表")
	sqlTable := `
		CREATE TABLE IF NOT EXISTS "user" (
		   "uid" INTEGER PRIMARY KEY AUTOINCREMENT,
		   "wx_id" VARCHAR(64) NULL,
		   "created" TIMESTAMP default (datetime('now', 'localtime'))  
		);
-- 		CREATE TABLE IF NOT EXISTS "userdeatail" (
-- 		   "uid" INT(10) NULL,
-- 		   "intro" TEXT NULL,
-- 		   "profile" TEXT NULL,
-- 		   PRIMARY KEY (uid)
-- 		);
		   `
	_,err = DB.Exec(sqlTable)
	public.Error(err)
}
