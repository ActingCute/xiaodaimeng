package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"xiaodaimeng/models"
	"xiaodaimeng/public"

	"time"
	"xiaodaimeng/controllor"
)

func init() {
	models.InitDB()
	//初始化日志文件
	file := "daimeng" + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("")
}

func main() {
	flag.Parse()
	u := url.URL{Scheme: "ws", Host: public.ConfigData.WsOrgin, Path: ""}
	public.Printf("connecting to ", u.String())

	c, _, connErr := websocket.DefaultDialer.Dial(u.String(), nil)
	if connErr != nil {
		public.Printf("dial:", connErr)
	}
	defer c.Close()
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				public.Error("read:", err)
				return
			}
			public.Debug("recv: ", string(message))
			controllor.Handle(message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	if !public.ConfigData.IsDebug {
		controllor.DoTimedTasks(controllor.SystemWxIdList.NeedNoticeUpdateList, "小呆萌开机了")
	}

	for {
		select {
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				public.Debug("write:", err)
				return
			}

		case rMsg := <-controllor.RWsMsg:
			err := c.WriteMessage(websocket.TextMessage, rMsg)
			//			print("\n进来关到了")
			if err != nil {
				public.Printf("write close:", err)
				return
			}
		}
	}

}
