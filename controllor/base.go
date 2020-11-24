package controllor

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
	"xiaodaimeng/public"
)

type SystemWxId struct {
	MaxAdminWxId         string   `json:"max_admin_wx_id"`         //最大的微信管理员
	NeedNoticeUpdateList []string `json:"need_notice_update_list"` //需要通知更新信息的id
}

var SystemWxIdList SystemWxId
var InWork map[string]bool //在工作中

//初始化一些json数据
func init() {
	InWork = make(map[string]bool)
	//初始化微信id
	wxIdPtr, _ := os.Open("data/wxIds.json")
	defer wxIdPtr.Close()
	decoder := json.NewDecoder(wxIdPtr)
	err := decoder.Decode(&SystemWxIdList)
	if err != nil {
		public.Printf("微信id解析失败，", err.Error())
	}

	//初始化表情
	emojiPtr, _ := os.Open("data/emoji.json")
	defer emojiPtr.Close()
	decoder = json.NewDecoder(emojiPtr)
	err = decoder.Decode(&emoji)
	if err != nil {
		public.Printf("表情解码失败，", err.Error())
	}
	//初始化一些问题和答案
	answerPtr, _ := os.Open("data/answer.json")
	defer answerPtr.Close()
	decoder = json.NewDecoder(answerPtr)
	err = decoder.Decode(&problemList)
	if err != nil {
		public.Printf("问题列表解码失败，", err.Error())
	}
}

//判断是不是被拉黑，是就不要傻傻的回复了，别人都拉黑了呢/或者是红包，不需要理会啦
func IsBlackMsg(msg Msg) bool {
	return strings.Index(msg.Content, "或系统消息") != -1 || strings.Index(msg.Content, "请在手机上查看") != -1
}

//判断是不是admin
func IsAdmin(msg Msg) bool {
	return isContains([]string{SystemWxIdList.MaxAdminWxId}, msg.MsgSender)
}

func isContains(father []string, son string) bool {
	has := false
	for _, key := range father {
		if key == son {
			has = true
			break
		}
	}
	return has
}

//md5
func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//是否在工作名单中
//判断是不是菜单函数
func IsInWork(wid string) bool {
	if _, ok := InWork[wid]; ok {
		return InWork[wid]
	}
	return false
}
