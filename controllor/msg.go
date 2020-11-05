package controllor

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// {"times":"2020-11-04 08-32-33","type":"文字","source":"群消息","wxid":"22925504714@chatroom","msgSender":"wxid_azmds1whb7r212","content":"啊"}

type Msg struct {
	Times     string `json:"times"`
	Type      string `json:"type"`
	Source    string `json:"source"`
	WxId      string `json:"wxid"`
	MsgSender string `json:"msgSender"`
	Content   string `json:"content"`
}

type RMsg struct {
	WxId    string `json:"m_wxid"`
	Content string `json:"m_Content"`
}

type XiaoDaiMeng struct {
	AnsNodeId    int    `json:"ans_node_id"`
	Answer       string `json:"answer"`
	FromUserName string `json:"from_user_name"`
}

type Emoji struct {
	Cn  string `json:"cn"`
	Zb  string `json:"zb"`
	img string `json:"img"`
}

type ProblemList struct {
	Problem string `json:"problem"`
	Answer  string `json:"answer"`
}

const (
	XiaoDaiMengName = "小呆萌"
	EncodingAESKey  = "cid2xSztPaRWLTVktma1tsO3rY9cD0d5SVRRW3AWgk3"
	OpenId          = "2X4xdKYO1kcByDWYZGqxgqqzi1zsQy"
	Token           = "puppet_donut_50d6bf5cbd5cdfa7"
)

const urlStr string = "https://openai.weixin.qq.com/openapi/message/" + OpenId

var RWsMsg chan []byte
var emoji []Emoji
var lock sync.Mutex
var oldEmojiKey = 0
var problemList []ProblemList

func init() {
	//初始化表情
	filePtr, _ := os.Open("data/emoji.json")
	defer filePtr.Close()
	decoder := json.NewDecoder(filePtr)
	err := decoder.Decode(&emoji)
	if err != nil {
		print("表情解码失败，", err.Error())
	}
	//初始化一些问题和答案
	//初始化表情
	filePtr1, _ := os.Open("data/answer.json")
	defer filePtr1.Close()
	decoder = json.NewDecoder(filePtr1)
	err = decoder.Decode(&problemList)
	if err != nil {
		print("问题列表解码失败，", err.Error())
	}
}

func getToken(msg Msg) (tokenString string, err error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["username"] = msg.WxId
	claims["msg"] = msg.Content
	token.Claims = claims
	tokenString, err = token.SignedString([]byte(EncodingAESKey))
	if err != nil {
		print(err)
		return
	}

	//print(tokenString)

	return
}

func httpPostForm(query string) (body []byte, err error) {
	resp, err := http.PostForm(urlStr,
		url.Values{"query": {query}})

	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	return
}

func Handle(bMsg []byte) {

	var msg Msg
	err := json.Unmarshal(bMsg, &msg)
	if err != nil {
		print("\nHandle Unmarshal error: ", err.Error())
		return
	}
	//自己发的，不用回复,已经拉黑的，不用回复
	if IsAdmin(msg) || IsBlackMsg(msg) {
		print("主动发出去的信息/红包 不回复")
		return
	}
	//判断是不是表情，是就用表情回复
	if msg.Type == "表情" {
		if len(emoji) > 1 {
			key := GenerateRangeNum(0, len(emoji)-1)
			if key == oldEmojiKey {

				if key == len(emoji)-1 {
					key = 0
				}

				if oldEmojiKey == 0 {
					key = 1
				} else {
					key = 0
				}
				oldEmojiKey = key
			}
			go func(index int) {
				SendMsg(msg.WxId, emoji[index].Cn)
			}(key)
		}
		return
	}
	//print("msg.Content:", msg.Content)
	token, err1 := getToken(msg)

	if err1 != nil {
		print("Handle getToken error: ", err1.Error())
		return
	}
	//print(token)
	msgBytes, err2 := httpPostForm(token)
	if err2 != nil {
		print("\nHandle httpPostForm error: ", err2.Error())
		return
	}
	print(string(msgBytes))
	var answer XiaoDaiMeng
	err = json.Unmarshal(msgBytes, &answer)
	if err != nil {
		print("\nHandle Unmarshal XiaoDaiMeng error: ", err.Error())
		return
	}

	if answer.AnsNodeId < 1 {
		//回答失败
		print("\nXiaoDaiMeng 回答失败了->", msg.Content, "<-\n")

		//查找问题列表有没有答案，有接直接用
		hasAnswer := false
		for _, a := range problemList {
			//print("strings.Index(msg.Content, a.Problem) ", strings.Index(msg.Content, a.Problem))
			if strings.Index(msg.Content, a.Problem) != -1 {
				answer.Answer = a.Answer
				hasAnswer = true
				break
			}
		}
		if !hasAnswer {
			return
		}
	}

	//print("\nCODE :", answer.AnsNodeId)

	//替换机器人的名字
	answer.Answer = strings.Replace(answer.Answer, "小微", XiaoDaiMengName, -1)
	//print("\n", answer.Answer)

	go func(a XiaoDaiMeng) {
		SendMsg(a.FromUserName, a.Answer)
	}(answer)
}

func SendMsg(wxId string, content string) {
	lock.Lock()
	time.Sleep(time.Microsecond * 8)
	var rMsg = RMsg{
		WxId:    wxId,
		Content: content,
	}
	bRMsg, err := json.Marshal(rMsg)
	if err != nil {
		print("\nSendMsg Marshal RMsg error: ", err.Error())
		return
	}
	//print("\n下来了")
	RWsMsg = make(chan []byte)
	RWsMsg <- bRMsg
	defer lock.Unlock()
}
