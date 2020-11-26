package models

import (
	"errors"
	"time"
	"xiaodaimeng/public"
)


type BlackList struct {
	Bid     int       `xorm:"int(20) pk not null autoincr 'bid'" json:"bid"`
	WxId    string    `xorm:"varchar(255) NOT NULL unique" json:"wx_id"`
	In      int       `xorm:"int(1) NULL "json:"in"` // 1 不自动聊天
	Created time.Time `xorm:"created" json:"created"`
	Updated time.Time `xorm:"updated" json:"updated"`
}

func SelectBlackList(list *[]BlackList) error {
	if DBOk {
		return DB.Table("black_list").Find(list)
	}
	return errors.New("数据库链接失败")
}

func GetBlackList(list *BlackList) (has bool,err error) {
	if DBOk {
		public.Debug("GetBlackList WxId : ",list.WxId)
		has,err  = DB.Table("black_list").Where("wx_id='wxid_u3q162gfuq8k22'").Get(list)
		public.Debug(list)
		if err != nil {
			public.Error(err)
		}
		return
	}
	return false,errors.New("数据库链接失败")
}

func insertBlackList(backList *BlackList) error {
	if DBOk {
		_, err := DB.Table("black_list").Insert(backList)
		if err != nil {
			public.Error(err)
		}
		return err
	}
	return errors.New("数据库链接失败")
}

func updateBlackList(backList *BlackList) error {
	if DBOk {
		_,err := DB.Table("black_list").Where("wx_id=?",backList.WxId).Update(backList)
		if err != nil {
			public.Error(err)
		}
		return  err
	}
	return errors.New("数据库链接失败")
}

func UpdateBlackList(backList *BlackList) error {
	if DBOk {
		bl := &BlackList{
			WxId:    backList.WxId,
			In:      backList.In,
		}
		//判断是否存在，不存在就插入
		public.Debug("UpdateBlackList :", backList.WxId)
		has,err := GetBlackList(bl)
		if err != nil {
			public.Error(err)
			return err
		}
		public.Debug("UpdateBlackList has : ",has)
		if !has {
			//不存在，需要插入
			err = insertBlackList(bl)
			if err != nil {
				return err
			}
		}
		//存在，更新
		err =  updateBlackList(bl)
		if err != nil {
			public.Error(err)
			return err
		}
		backList = bl
		return nil
	}
	return errors.New("数据库链接失败")
}
