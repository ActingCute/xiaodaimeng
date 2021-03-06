package controllor

const (
	SuccessCode = 6666
	//下面的是跟微信助手约定好的
	WMSendTextMessage  = 5
	WMSendImageMessage = 8

	HEART_BEAT             = 5005
	RECV_TXT_MSG           = 1
	RECV_PIC_MSG           = 3
	USER_LIST              = 5000
	GET_USER_LIST_SUCCSESS = 5001
	GET_USER_LIST_FAIL     = 5002
	TXT_MSG                = 555
	PIC_MSG                = 500
	AT_MSG                 = 550
	CHATROOM_MEMBER        = 5010
	CHATROOM_MEMBER_NICK   = 5020
	PERSONAL_INFO          = 6500
	DEBUG_SWITCH           = 6000
	PERSONAL_DETAIL        = 6550
)
const (
	FailText    = "小呆萌开了小差嘤嘤嘤，完成不了工作"
	NotDoneText = "还在开发中呢~"
	HasDrawText = "啊！你今天已经抽了签啦~\n回复8查看解签"
	NotDrawText = "啊！你今天还没抽签喔~回复7抽签"
	XiaoDaiMengCried = "嘤嘤嘤,你就这么讨厌小呆萌嘛~"
	XiaoDaiMengLose = "桀桀桀,小灰已上线，可以发 小呆萌 叫小呆萌回来"
	XiaoDaiMengStay = "小呆萌一直在喔~"
	XiaoDaiMengSleep = "小呆萌已经睡着了，可以发 小呆萌 叫小呆萌醒来"
	XiaoDaiMengCome = "嘿嘿，小呆萌回来啦~"
)

const Self = "self"