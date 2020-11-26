package controllor

import (
	"xiaodaimeng/public"
)

//聊天黑名单
var Blacklist = map[string]bool{}


//关闭自动聊天
func OffCatch(msg Msg)  {
	public.Debug("关闭自动聊天",IsInBlacklist(msg))
	//判断是否已经关闭
	sid := GetSender(msg)
	rid := GetReceiver(msg)
	if IsInBlacklist(msg) {
		SendMsg(rid, XiaoDaiMengCried, TXT_MSG)
		return
	}
	//bl := models.BlackList{
	//	WxId:    sid,
	//	In:      1}
	//err := models.UpdateBlackList(&bl)
	//if err != nil {
	//	SendMsg(rid, FailText, TXT_MSG)
	//	public.Error("OffCatch:",err)
	//	return
	//}
	Blacklist[sid] = true
	SendMsg(rid, XiaoDaiMengLose, TXT_MSG)
}

//开启自动聊天
func OnCatch(msg Msg)  {
	public.Debug("开启自动聊天",IsInBlacklist(msg))
	//判断是否已经关闭
	sid := GetSender(msg)
	rid := GetReceiver(msg)
	if !IsInBlacklist(msg) {
		SendMsg(rid, XiaoDaiMengStay, TXT_MSG)
		return
	}
	//bl := models.BlackList{
	//	WxId:    sid,
	//	In:     2}
	//err := models.UpdateBlackList(&bl)
	//if err != nil {
	//	SendMsg(rid, FailText, TXT_MSG)
	//	public.Error("OnCatch:",err)
	//	return
	//}
	Blacklist[sid] = false
	public.Debug("IsInBlacklist(msg) - ",IsInBlacklist(msg),Blacklist[sid],sid)
	SendMsg(rid, XiaoDaiMengCome, TXT_MSG)
}

//判断是不是在黑名单
func IsInBlacklist(msg Msg) bool {
	sid := GetSender(msg)
	public.Debug("IsInBlacklist sid:",sid)
	if _, ok := Blacklist[sid]; ok {
		return true
	}
	return false
}