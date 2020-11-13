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
	Wid     int       `xorm:"int(20) pk not null autoincr 'wid'" json:"wid"`
	WxId    string    `xorm:"varchar(50) NOT NULL" json:"wx_id"`
	Type    string    `xorm:"varchar(50) NOT NULL"json:"type"`
	Msg     string    `xorm:"varchar(50) NULL "json:"msg"`
	Other   string    `xorm:"varchar(50) NULL "json:"other"`
	Created time.Time `xorm:"created" json:"created"`
}

func SelectWork(work *Work) error {
	if DBOk {
		DB.Table("work").Where("wx_id=? and type = ? and msg = ?", work.WxId, work.Type, work.Msg).Get(work)
		return nil
	}
	return errors.New("数据库链接失败")
}

func InsertWork(work *Work) error {
	_, err := DB.Table("work").Insert(work)
	if err != nil {
		return err
	}
	return SelectWork(work)
}
