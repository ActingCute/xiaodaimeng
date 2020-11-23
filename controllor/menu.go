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

type XDM struct {
	Menu   []string `json:"menu"`
	About  string   `json:"about"`
	Update Update   `json:"update"`
}

type Update struct {
	Version int          `json:"version"` //当前
	Info    []UpdateInfo `json:"info"`
}

type UpdateInfo struct {
	Version int    `json:"version"`
	Msg     string `json:"msg"`
}

var Menu = new(XDM)
var menuFunc = map[string]func(msg Msg){}
var MenuText = ""

func init() {
	//获取菜单
	menuPtr, _ := os.Open("data/menu.json")
	defer menuPtr.Close()
	decoder := json.NewDecoder(menuPtr)
	err := decoder.Decode(&Menu)
	if err != nil {
		public.Printf("菜单文件解码失败，", err.Error())
		return
	}
	public.Debug("about : ", Menu.About)
	for i, f := range Menu.Menu {
		funcNames := strings.Split(f, "/")
		ff := func(msg Msg) {}
		//"帮助/h",
		//	"成语接龙/cyjl",
		//	"自定义问题答案/pa",
		//	"日记本/diary",
		//	"待办事项/todo"

		if len(funcNames) == 2 {
			public.Debug(funcNames[0])
			switch funcNames[0] {
			case "帮助":
				ff = Help
				break
			case "关于":
				ff = About
				break
			case "成语接龙":
				ff = Cyjl
				break
			case "自定义问题答案":
				ff = Pa
				break
			case "日记":
				ff = Diary
				break
			case "待办事项":
				ff = WTodo
				break
			case "更新信息":
				ff = GetUpdateInfo
				break
			case "抽签":
				ff = Draw
			case "解签":
				ff = UnDraw
				break
			}
			key := strconv.Itoa(i)
			menuFunc[md5V(strings.ToUpper(funcNames[0]))] = ff
			menuFunc[md5V(strings.ToUpper(funcNames[1]))] = ff
			menuFunc[md5V(key)] = ff
			MenuText += key + ". " + f + "\n"
		}

	}

	if MenuText != "" {
		MenuText += "\n直接回复序号或文字获取功能"
	}

}

//菜单函数

//判断是不是菜单函数
func IsMenuFunc(ff string) func(msg Msg) {
	ff = md5V(strings.ToUpper(ff))
	if _, ok := menuFunc[ff]; ok {
		return menuFunc[ff]
	}
	return nil
}

//帮助
func Help(msg Msg) {
	if MenuText != "" {
		SendMsg(msg.Sender, MenuText, TXT_MSG)
	} else {
		//没有菜单
		GetAnswer(msg)
	}
}

//关于
func About(msg Msg) {
	public.Debug("About")
	SendMsg(msg.Sender, Menu.About, TXT_MSG)
}

//成语接龙
func Cyjl(msg Msg) {
	public.Debug("Cyjl")
	msg.Content = "成语接龙"
	GetAnswer(msg)
}

//自定义问题答案
func Pa(msg Msg) {
	public.Debug("Pa")
	SendMsg(msg.Sender, NotDoneText, TXT_MSG)
}

//日记
func Diary(msg Msg) {
	public.Debug("Diary")
	SendMsg(msg.Sender, NotDoneText, TXT_MSG)
}

//待办事项
func WTodo(msg Msg) {
	public.Debug("Todo")
	SendMsg(msg.Sender, NotDoneText, TXT_MSG)
}

//获取更新信息
func GetUpdateInfo(msg Msg) {
	public.Debug("GetUpdateInfo")

	updateInfo := Menu.Update.Info[0].Msg

	for _, info := range Menu.Update.Info {
		if info.Version == Menu.Update.Version {
			updateInfo = info.Msg
			break
		}
	}

	updateInfo = "当前版本：" + strconv.Itoa(Menu.Update.Version) + "\n" + updateInfo

	SendMsg(msg.Sender, updateInfo, TXT_MSG)
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
	public.Error(err)

	if err != nil {
		public.Error(err)
		SendMsg(msg.Sender, FailText, TXT_MSG)
		return
	}
	if work.Wid > 0 {
		//您今天已经抽了签
		SendMsg(msg.Sender, HasDrawText, TXT_MSG)
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

	SendMsg(msg.Sender, linQian.Number, TXT_MSG)

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
		SendMsg(msg.Sender, FailText, TXT_MSG)
		return
	}
	if work.Wid < 1 {
		//您今天还没抽签喔
		SendMsg(msg.Sender, NotDrawText, TXT_MSG)
		return
	}
	linQianInfo := strings.Split(work.Other, "/")

	if len(linQianInfo) < 2 {
		SendMsg(msg.Sender, FailText, TXT_MSG)
		return
	}

	linQianType := linQianInfo[0]             //签类型
	key, err1 := strconv.Atoi(linQianInfo[1]) //第几签
	if err1 != nil {
		public.Error(err1)
		SendMsg(msg.Sender, FailText, TXT_MSG)
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
			SendMsg(msg.Sender, lcStr, TXT_MSG)
			lcStr = ""
			isSend = true
		}
		lcStr += lq + "\n"
		if i == len(linQian.Content)-1 {
			if !isSend {
				SendMsg(msg.Sender, lcStr, TXT_MSG)
			}
			InWork[msg.Sender] = false //移除工作名单
		}
	}
}
