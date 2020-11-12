package controllor

import (
	"encoding/json"
	"os"
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
