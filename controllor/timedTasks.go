package controllor

import (
	"encoding/json"
	"github.com/robfig/cron"
	"os"
	"sync"
)

var FriendsWeiXinId  = []string{"wxid_dppjrkktwdfd22","wxid_ch5ql4b1uh2v22","wxid_u3q162gfuq8k22","wxid_vw6bngu8c4o721","wxid_azmds1whb7r212","wxid_humu0ux5uz5622"}
//var FriendsWeiXinId = []string{"wxid_u3q162gfuq8k22", "wxid_azmds1whb7r212"}
var AdminWeiXinId = []string{"qq190025254"}

const (
	sayGoodNightText = "睡觉啦，晚安~"
	weiBoSignText    = "新浪要签到啦~"
	drinkWaterText   = "要喝水啦~"
	curriculumText   = "接下来的课程是： "
)

type CurriculumContent struct {
	Name      string `json:"name"`
	StartTime string `json:"start_time"`
	Number    int    `json:"number"`
}

var Curriculum []CurriculumContent

func init() {
	print("\ninit goodNight\n")
	c := cron.New()
	//晚安
	c.AddFunc("00 00 00 * * ?", func() {
		doTimedTasks(FriendsWeiXinId, sayGoodNightText)
	})
	//新浪签到
	c.AddFunc("00 00 22 * * ?", func() {
		doTimedTasks(AdminWeiXinId, weiBoSignText)
	})
	//要喝水啦
	c.AddFunc("0 0 08,09,10,11,12,13,14,15,16,17,18,19,20,21,22 * * ?", func() {
		doTimedTasks(AdminWeiXinId, drinkWaterText)
	})
	//课程
	filePtr, _ := os.Open("data/curriculum.json")
	defer filePtr.Close()
	decoder := json.NewDecoder(filePtr)
	err := decoder.Decode(&Curriculum)
	if err != nil {
		print("课程解码失败，", err)
	} else {
		var wg sync.WaitGroup
		for _, curriculum := range Curriculum {
			wg.Add(1)
			go func(curr CurriculumContent) {
				c.AddFunc(curr.StartTime, func() {
					doTimedTasks(AdminWeiXinId, curriculumText+curr.Name)
				})
				wg.Done()
			}(curriculum)
		}
		wg.Wait()
	}

	c.Start()
}

//说晚安
func doTimedTasks(wxIds []string, msg string) {
	for _, wxId := range wxIds {
		var rMsg = RMsg{
			WxId:    wxId,
			Content: msg,
		}
		bRMsg, err := json.Marshal(rMsg)
		if err != nil {
			print("\ndoTimedTasks Marshal RMsg error: ", err.Error())
			continue
		}
		RWsMsg = make(chan []byte)
		RWsMsg <- bRMsg
	}
}
