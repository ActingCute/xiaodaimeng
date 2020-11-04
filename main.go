package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"

	"time"
	"xiaodaimeng/controllor"
)

var addr = flag.String("addr", "localhost:9100", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "ws", Host: *addr, Path: ""}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
 
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			//log.Printf("recv: %s", message)
			controllor.Handle(message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case rMsg := <-controllor.RWsMsg:
			err := c.WriteMessage(websocket.TextMessage, rMsg)
//			print("\n进来关到了")
			if err != nil {
				log.Println("write close:", err)
				return
			}
		}
	}
}
