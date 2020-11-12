package models

import (
	"errors"
	"time"
	"xiaodaimeng/public"
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
	Created time.Time `json:"created"`
}

func SelectWork(work *Work) error {
	if DBOk {
		rows, err := DB.Query("SELECT * FROM work ")
			//where `wx_id` = ? and `type` = ? and `msg` = ?", work.WxId, work.Type, work.Msg)
		for rows.Next() {
			err = rows.Scan(work)
			if err != nil {
				public.Error(err)
			}
		}
		return err
	}
	return errors.New("数据库链接失败")
}

func InsertWork(work *Work) error {
	stmt, err := DB.Prepare("INSERT INTO work(wx_id, type, msg) values(?,?,?)")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(work.WxId, work.Type, work.Msg)

	if err != nil {
		return err
	}

	return SelectWork(work)
}
