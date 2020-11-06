package controllor

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
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

const notDoneText = "还在开发中呢~"

func init() {
	//获取菜单
	menuPtr, _ := os.Open("data/menu.json")
	defer menuPtr.Close()
	decoder := json.NewDecoder(menuPtr)
	err := decoder.Decode(&Menu)
	if err != nil {
		Printf("菜单文件解码失败，", err.Error())
	}
	Debug("about : ", Menu.About)
	for i, f := range Menu.Menu {
		funcNames := strings.Split(f, "/")
		ff := func(msg Msg) {}
		//"帮助/h",
		//	"成语接龙/cyjl",
		//	"自定义问题答案/pa",
		//	"日记本/diary",
		//	"待办事项/todo"

		if len(funcNames) == 2 {
			ConfigData.IsDebug = true
			Debug(funcNames[0])
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

func Help(msg Msg) {
	if MenuText != "" {
		SendMsg(msg.WxId, MenuText)
	} else {
		//没有菜单
		GetAnswer(msg)
	}
}

func About(msg Msg) {
	Debug("About")
	SendMsg(msg.WxId, Menu.About)
}

func Cyjl(msg Msg) {
	Debug("Cyjl")
	GetAnswer(msg)
}

func Pa(msg Msg) {
	Debug("Pa")
	SendMsg(msg.WxId, notDoneText)
}

func Diary(msg Msg) {
	Debug("Diary")
	SendMsg(msg.WxId, notDoneText)
}

func WTodo(msg Msg) {
	Debug("Todo")
	SendMsg(msg.WxId, notDoneText)
}
