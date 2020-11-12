package models

import (
	"errors"
	"time"
)

//"wid" INTEGER PRIMARY KEY AUTOINCREMENT,
//"wx_id" VARCHAR(64) NOT NULL,
//"type" TEXT NOT NULL,
//"msg" TEXT NOT NULL,
//"created" TIMESTAMP default (datetime('now', 'localtime'))

type Work struct {
	Wid     int       `json:"wid"`
	WxId    string    `json:"wx_id"`
	Type    string    `json:"type"`
	Msg     string    `json:"msg"`
	Other   string    `json:"other"`
	Created time.Time `json:"created"`
}

func SelectWork(work *Work) error {
	if DBOk {
		DB.Table("work").Where("wx_id=? and type = ? and msg = ?", work.WxId, work.Type, work.Msg).Get(work)
		return nil
	}
	return errors.New("数据库链接失败")
}

func InsertWork(work *Work) error {
	err := DB.Table("work").Insert()
	return SelectWork(work)
}
