package controllor

import (
	"encoding/json"
	"github.com/robfig/cron"
	"os"
	"strconv"
	"sync"
	"time"
	"xiaodaimeng/public"
)

const (
	curriculumText = "接下来的课程是： "
)

type CurriculumContent struct {
	Name      string `json:"name"`
	StartTime string `json:"start_time"`
	Number    int    `json:"number"`
}

type TimedTasks struct {
	Time string   `json:"time"`
	Msg  string   `json:"msg"`
	WxId []string `json:"wx_id"`
	Type string   `json:"type"`
}

var Curriculum []CurriculumContent
var TimedTasksList []TimedTasks

func init() {

	c := cron.New()

	//一些定时任务
	timedTasksPtr, _ := os.Open("data/timedTasks.json")
	defer timedTasksPtr.Close()
	decoder := json.NewDecoder(timedTasksPtr)
	err := decoder.Decode(&TimedTasksList)

	if err != nil {
		public.Printf("定时任务解析失败，", err)
	} else {
		var wg sync.WaitGroup
		for _, tt := range TimedTasksList {
			public.Debug(tt)
			wg.Add(1)
			go func(curr TimedTasks) {
				c.AddFunc(curr.Time, func() {
					DoTimedTasksWithType(curr.WxId, curr.Msg, curr.Type)
				})
				wg.Done()
			}(tt)
		}
		wg.Wait()
	}

	//课程
	filePtr, _ := os.Open("data/curriculum.json")
	defer filePtr.Close()
	decoder = json.NewDecoder(filePtr)
	err = decoder.Decode(&Curriculum)
	if err != nil {
		public.Error("课程解码失败，", err)
	} else {
		var wg sync.WaitGroup
		for _, curriculum := range Curriculum {
			wg.Add(1)
			go func(curr CurriculumContent) {
				c.AddFunc(curr.StartTime, func() {
					DoTimedTasks([]string{SystemWxIdList.MaxAdminWxId}, curriculumText+curr.Name)
				})
				wg.Done()
			}(curriculum)
		}
		wg.Wait()
	}
	c.Start()
}

//处理定时任务
func DoTimedTasks(wxIds []string, msg string) {
	for _, wxId := range wxIds {
		go SendMsg(wxId, msg, TXT_MSG)
	}
}

func DoTimedTasksWithType(wxIds []string, msg, tp string) {
	t := time.Now()
	w := int(t.Weekday())

	public.Debug("w--",w)
	for _, wxId := range wxIds {
		if tp == "good_night" {
			picName := public.GetCurrentDirectory() + "/static/img/goodnight/" + strconv.Itoa(w) + ".jpg"
			go SendMsg(wxId, picName, PIC_MSG)
		}else{
			go SendMsg(wxId, msg, TXT_MSG)
		}
	}
}
