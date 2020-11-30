package models

import (
	"errors"
	"time"
	"xiaodaimeng/public"
)


type BlackList struct {
	Bid     int       `xorm:"int(20) pk not null autoincr 'bid'" json:"bid"`
	WxId    string    `xorm:"varchar(255) NOT NULL unique 'wx_id'" json:"wx_id"`
	Open    int       `xorm:"int(1) default 2 "json:"open"` // 2 不自动聊天
	Created time.Time `xorm:"created" json:"created"`
	Updated time.Time `xorm:"updated" json:"updated"`
}

var BlacklistMap = map[string]bool{}

func init()  {
	BlacklistMap = make(map[string]bool)
}

func SelectBlackList(list *[]BlackList) error {
	if DBOk {
		err := DB.Table(TableBlackList).Find(list)
		if err != nil {
			public.Error(err)
		}
		return err
	}
	return errors.New("数据库链接失败")
}

func GetBlackList(list *BlackList) (has bool,err error) {
	if DBOk {
		//直接用get不知道为啥报唯一错
        ls := make([]BlackList,0)
		err = DB.Table(TableBlackList).Where("wx_id=?", list.WxId).Find(&ls)
		has = len(ls) > 0
		if err != nil {
			public.Error(err)
			return
		}
		if len(ls) > 0 {
			*list = ls[0]
		}
		return
	}
	return false,errors.New("数据库链接失败")
}

func insertBlackList(backList *BlackList) error {
	if DBOk {
		_, err := DB.Table(TableBlackList).Insert(backList)
		if err != nil {
			public.Error(err)
		}
		return err
	}
	return errors.New("数据库链接失败")
}

func updateBlackList(backList *BlackList) error {
	if DBOk {
		public.Debug("backList    === ",backList)
		_,err := DB.Table(TableBlackList).Where("wx_id=?",backList.WxId).Omit("in").Update(backList)
		if err != nil {
			public.Error(err)
		}
		return  err
	}
	return errors.New("数据库链接失败")
}

func UpdateBlackList(backList BlackList) error {
	if DBOk {
		//判断是否存在，不存在就插入
		public.Debug("UpdateBlackList :", backList.WxId)
		checkBackList := BlackList{
			Bid:     backList.Bid,
			WxId:    backList.WxId,
			Open:      backList.Open,
			Created: backList.Created,
			Updated: backList.Updated,
		}
		has,err := GetBlackList(&checkBackList)
		if err != nil {
			public.Error(err)
			return err
		}
		public.Debug("UpdateBlackList has : ",has,backList)
		if !has {
			//不存在，需要插入
			public.Debug("不存在，需要插入")
			return insertBlackList(&backList)
		}
		//存在，更新
		public.Debug("存在，更新")

		err =  updateBlackList(&backList)
		return err
	}
	return errors.New("数据库链接失败")
}

//加载列表缓存
func InitBlacklist()  {
	list := make([]BlackList,0)
	err := SelectBlackList(&list)
	if err == nil {
		for _,bl := range list {
			BlacklistMap[bl.WxId] = bl.Open == 1
			public.Debug(bl.WxId," : ",BlacklistMap[bl.WxId])
		}
	}else {
		public.Error(err)
	}
}