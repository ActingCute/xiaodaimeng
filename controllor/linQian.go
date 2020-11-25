package controllor

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"
	"xiaodaimeng/models"
	"xiaodaimeng/public"
)

type Lucky struct {
	Key     string   `json:"key"`
	Number  string   `json:"number"`
	Content []string `json:"content"`
}

type LuckyData struct {
	GuanYin []Lucky `json:"guan_yin"`
	YueLao  []Lucky `json:"yue_lao"`
}

var LuckyDataList = new(LuckyData)

//抽签

func init() {
	//获取灵签数据

	//观音
	guanYinLingQianPtr, _ := os.Open("data/guanyinlingqian.json")
	defer guanYinLingQianPtr.Close()
	decoder := json.NewDecoder(guanYinLingQianPtr)
	err := decoder.Decode(&LuckyDataList.GuanYin)
	if err != nil {
		public.Error("观音签数据解码失败，", err.Error())
	}

	//月老
	yueLaoLingQianPtr, _ := os.Open("data/yuelaolingqian.json")
	defer yueLaoLingQianPtr.Close()
	decoder = json.NewDecoder(yueLaoLingQianPtr)
	err = decoder.Decode(&LuckyDataList.YueLao)
	if err != nil {
		public.Error("月老签数据解码失败，", err.Error())
	}
}


//抽签
func Draw(msg Msg) {
	public.Debug("Draw")
	//判断今天是否已经抽签
	work := models.Work{
		WxId: msg.Sender,
		Type: "draw",
		Msg:  time.Now().Format("2006-01-02"),
	}
	err := models.SelectWork(&work)

	if err != nil {
		public.Error(err)
		SendMsg(GetReceiver(msg), FailText, TXT_MSG)
		return
	}
	if work.Wid > 0 {
		//您今天已经抽了签
		SendMsg(GetReceiver(msg), HasDrawText, TXT_MSG)
		return
	}
	//随机签的类型
	key := public.GenerateRangeNum(0, 1)
	linQianList := LuckyDataList.GuanYin
	lenLinQian := len(LuckyDataList.GuanYin)
	linQianType := "GuanYin"
	if key == 1 {
		linQianList = LuckyDataList.YueLao
		lenLinQian = len(LuckyDataList.YueLao)
		linQianType = "YueLao"
	}
	key = public.GenerateRangeNum(0, lenLinQian)
	linQian := linQianList[key]

	SendMsg(GetReceiver(msg), linQian.Number, TXT_MSG)

	//插入数据库
	work.Other = linQianType + "/" + linQian.Key
	go models.InsertWork(&work)
}

//解签
func UnDraw(msg Msg) {
	public.Debug("UnDraw")
	//判断今天是否已经抽签
	work := models.Work{
		WxId: msg.Sender,
		Type: "draw",
		Msg:  time.Now().Format("2006-01-02"),
	}
	err := models.SelectWork(&work)
	if err != nil {
		public.Error(err)
		SendMsg(GetReceiver(msg), FailText, TXT_MSG)
		return
	}
	if work.Wid < 1 {
		//您今天还没抽签喔
		SendMsg(GetReceiver(msg), NotDrawText, TXT_MSG)
		return
	}
	linQianInfo := strings.Split(work.Other, "/")

	if len(linQianInfo) < 2 {
		SendMsg(GetReceiver(msg), FailText, TXT_MSG)
		return
	}

	linQianType := linQianInfo[0]             //签类型
	key, err1 := strconv.Atoi(linQianInfo[1]) //第几签
	if err1 != nil {
		public.Error(err1)
		SendMsg(GetReceiver(msg), FailText, TXT_MSG)
		return
	}
	InWork[msg.Sender] = true //加入工作名单
	linQianList := LuckyDataList.GuanYin
	if linQianType == "YueLao" {
		linQianList = LuckyDataList.YueLao
	}

	linQian := linQianList[key]

	lcStr := ""
	for i, lq := range linQian.Content {
		isSend := false
		if len(lcStr+lq) > 2000 {
			SendMsg(GetReceiver(msg), lcStr, TXT_MSG)
			lcStr = ""
			isSend = true
		}
		lcStr += lq + "\n"
		if i == len(linQian.Content)-1 {
			if !isSend {
				SendMsg(GetReceiver(msg), lcStr, TXT_MSG)
			}
			InWork[msg.Sender] = false //移除工作名单
		}
	}
}
