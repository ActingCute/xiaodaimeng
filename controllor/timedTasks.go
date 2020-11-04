package controllor

import (
	"encoding/json"
	"github.com/robfig/cron"
	"os"
	"sync"
)



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
	c := cron.New()
	//晚安
	c.AddFunc("00 00 00 * * ?", func() {
		doTimedTasks(FriendsWeiXinId, sayGoodNightText)
	})
	//新浪签到
	c.AddFunc("00 00 22 * * ?", func() {
		doTimedTasks(MaxAdminIds, weiBoSignText)
	})
	//要喝水啦
	c.AddFunc("0 0 08,09,10,11,12,13,14,15,16,17,18,19,20,21,22 * * ?", func() {
		doTimedTasks(MaxAdminIds, drinkWaterText)
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

//处理定时任务
func doTimedTasks(wxIds []string, msg string) {
	for _, wxId := range wxIds {
		go SendMsg(wxId, msg)
	}
}
