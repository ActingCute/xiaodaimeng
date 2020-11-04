package controllor

import (
	"math/rand"
	"time"
)

var FriendsWeiXinId = []string{"wxid_dppjrkktwdfd22", "wxid_ch5ql4b1uh2v22", "wxid_u3q162gfuq8k22", "wxid_vw6bngu8c4o721", "wxid_azmds1whb7r212", "wxid_humu0ux5uz5622"}

//var FriendsWeiXinId = []string{"wxid_u3q162gfuq8k22", "wxid_azmds1whb7r212"}
var AdminWeiXinId = []string{"Conan444444164", "wxid_azmds1whb7r212"}

const MaxAdminId = "wxid_u3q162gfuq8k22"

var MaxAdminIds = [] string{MaxAdminId, "qq190025254"}

//两个数之间的随机数
func GenerateRangeNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max-min) + min
	return randNum
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
