package controllor

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"os"
	"strings"
	"time"
	"xiaodaimeng/public"
)

var FriendsWeiXinId = []string{"wxid_ppiqsunhwmho22", "wxid_dppjrkktwdfd22", "wxid_ch5ql4b1uh2v22", "wxid_u3q162gfuq8k22", "wxid_vw6bngu8c4o721", "wxid_azmds1whb7r212", "wxid_humu0ux5uz5622"}

//var FriendsWeiXinId = []string{"wxid_u3q162gfuq8k22", "wxid_azmds1whb7r212"}
var AdminWeiXinId = []string{"Conan444444164", "wxid_azmds1whb7r212"}

const MaxAdminId = "wxid_u3q162gfuq8k22"

var NeedSendUpdateList = []string{"wxid_ppiqsunhwmho22", MaxAdminId} //需要通知更新的名单
var DrinkWaterList = []string{"wxid_ppiqsunhwmho22", MaxAdminId}     //需要提醒喝水的列表
var MaxAdminIds = []string{MaxAdminId}

//初始化一些json数据
func init() {
	//初始化表情
	emojiPtr, _ := os.Open("data/emoji.json")
	defer emojiPtr.Close()
	decoder := json.NewDecoder(emojiPtr)
	err := decoder.Decode(&emoji)
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

//两个数之间的随机数
func GenerateRangeNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max-min) + min
	return randNum
}

//判断是不是被拉黑，是就不要傻傻的回复了，别人都拉黑了呢/或者是红包，不需要理会啦
func IsBlackMsg(msg Msg) bool {
	return strings.Index(msg.Content, "或系统消息") != -1
}

//判断是不是admin
func IsAdmin(msg Msg) bool {
	return isContains(AdminWeiXinId, msg.MsgSender) || isContains(AdminWeiXinId, msg.MsgSender)
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
