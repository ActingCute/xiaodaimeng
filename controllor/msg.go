package controllor

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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

const (
	XiaoDaiMengName = "小呆萌"
	EncodingAESKey  = "cid2xSztPaRWLTVktma1tsO3rY9cD0d5SVRRW3AWgk3"
	OpenId          = "2X4xdKYO1kcByDWYZGqxgqqzi1zsQy"
	Token           = "puppet_donut_50d6bf5cbd5cdfa7"
)


const urlStr string = "https://openai.weixin.qq.com/openapi/message/" + OpenId

var RWsMsg chan []byte

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

	if answer.AnsNodeId > 0 {
		//回答失败
		print("\nXiaoDaiMeng 回答失败\n",answer.Answer)
		return
	}

	//替换机器人的名字
	answer.Answer = strings.Replace(answer.Answer, "小微", XiaoDaiMengName, -1)
	print("\n", answer.Answer)

	var rMsg = RMsg{
		WxId:    answer.FromUserName,
		Content: answer.Answer,
	}
	bRMsg, err := json.Marshal(rMsg)
	if err != nil {
		print("\nHandle Marshal RMsg error: ", err.Error())
		return
	}
	//print("\n下来了")
	RWsMsg = make(chan []byte)
	RWsMsg <- bRMsg
}
