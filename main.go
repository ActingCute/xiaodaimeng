package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"

	"time"
	"xiaodaimeng/controllor"
)

func init() {

	//初始化日志文件
	file := "daimeng" + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("")

	//初始化配置
	configPtr, _ := os.Open("data/config.json")
	defer configPtr.Close()
	decoder := json.NewDecoder(configPtr)
	err = decoder.Decode(&controllor.ConfigData)
	if err != nil {
		controllor.Printf("配置文件解码失败，", err.Error())
		//给一些默认的值
		controllor.ConfigData.IsDebug = false
	}
}

func main() {
	flag.Parse()
	u := url.URL{Scheme: "ws", Host: controllor.ConfigData.WsOrgin, Path: ""}
	controllor.Printf("connecting to ", u.String())

	c, _, connErr := websocket.DefaultDialer.Dial(u.String(), nil)
	if connErr != nil {
		controllor.Printf("dial:", connErr)
	}
	defer c.Close()
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				controllor.Printf("read:", err)
				return
			}
			controllor.Debug("recv: ", string(message))
			controllor.Handle(message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	if !controllor.ConfigData.IsDebug {
		controllor.DoTimedTasks(controllor.NeedSendUpdateList, "小呆萌开机了")
	}

	for {
		select {
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				controllor.Debug("write:", err)
				return
			}

		case rMsg := <-controllor.RWsMsg:
			err := c.WriteMessage(websocket.TextMessage, rMsg)
			//			print("\n进来关到了")
			if err != nil {
				controllor.Printf("write close:", err)
				return
			}
		}
	}

}
