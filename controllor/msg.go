package controllor

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
	"xiaodaimeng/public"
)

// {"times":"2020-11-04 08-32-33","type":"文字","source":"群消息","wxid":"22925504714@chatroom","msgSender":"wxid_azmds1whb7r212","content":"啊"}

//{
//id:getid(),
//type:CHATROOM_MEMBER_NICK,
//content:'5325308046@chatroom',//chatroom id 23023281066@chatroom  17339716569@chatroom
////5325308046@chatroom
////5629903523@chatroom
//wxid:'ROOT'
//  }

//{"content":"","id":"","sender":"ROOT","srvid":1,"time":"2020-11-23 21:13:51","type":5005}

type Msg struct {
	Id        string `json:"id"`
	Time      string `json:"time"`
	Type      int    `json:"type"`
	Sender    string `json:"sender"`
	WxId      string `json:"wxid"`
	MsgSender string `json:"msgSender"`
	Content   string `json:"content"`
	SrvId     int    `json:"srvid"`
}

type XiaoDaiMeng struct {
	AnsNodeId    int    `json:"ans_node_id"`
	Answer       string `json:"answer"`
	FromUserName string `json:"from_user_name"`
	AnswerType   string `json:"answer_type"`
	MoreInfo     MoreInfo `json:"more_info"`
}

type MoreInfo struct {
	NewsAnsDetail string `json:"news_ans_detail"`
}

type MoreInfoData struct {
	MoreInfo NewsData `json:"data"`
}

type NewsData struct {
	NewsDocs []NewsDocs `json:"docs"`
}

type NewsDocs struct {
	AbsL 		string `json:"abs_l"`
	AbsM 		string `json:"abs_m"`
	AbsS 		string `json:"abs_s"`
	Cate1 		string `json:"cate1"`
	Cate2 		string `json:"cate2"`
	DocId 		string `json:"docid"`
	Pubtime 	string `json:"pubtime"`
	Shortcut	string `json:"shortcut"`
	Srcfrom		string `json:"srcfrom"`
	Title 		string `json:"title"`
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

	NewsAnswer = "news"
)

const urlStr string = "https://openai.weixin.qq.com/openapi/message/" + OpenId

var RWsMsg chan []byte
var emoji []Emoji
var lock sync.Mutex
var oldEmojiKey = 0
var problemList []ProblemList

func getTimeId() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func getToken(msg Msg) (tokenString string, err error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["username"] = msg.Sender
	claims["msg"] = msg.Content
	token.Claims = claims
	tokenString, err = token.SignedString([]byte(EncodingAESKey))
	if err != nil {
		public.Printf(err.Error())
		return
	}

	//Printf(tokenString)

	return
}

func httpPostForm(query string) (body []byte, err error) {
	resp, err := http.PostForm(urlStr,
		url.Values{"query": {query}})

	if err != nil {
		public.Error(err.Error())
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
		public.Error("\nHandle Unmarshal error: ", err.Error())
		return
	}
	//心跳
	if msg.Type == HEART_BEAT {
		return
	}
	//非文字消息
	if msg.Type != RECV_TXT_MSG {
		return
	}
	//自己发的，不用回复,已经拉黑的，不用回复
	if IsAdmin(msg) || IsBlackMsg(msg) || IsInWork(msg.Sender) {
		public.Printf("在工作名单中/红包 不回复")
		return
	}
	//判断是不是菜单函数
	if ff := IsMenuFunc(msg.Content); ff != nil {
		public.Debug("是菜单函数")
		ff(msg)
		return
	}
	//判断是不是表情，是就用表情回复
	//if msg.Type == "表情" {
	//	if len(emoji) > 1 {
	//		key := public.GenerateRangeNum(0, len(emoji)-1)
	//		if key == oldEmojiKey {
	//
	//			if key == len(emoji)-1 {
	//				key = 0
	//			}
	//
	//			if oldEmojiKey == 0 {
	//				key = 1
	//			} else {
	//				key = 0
	//			}
	//			oldEmojiKey = key
	//		}
	//		go func(index int) {
	//			SendMsg(msg.WxId, emoji[index].Cn)
	//		}(key)
	//	}
	//	return
	//}
	//Printf("msg.Content:", msg.Content)
	//机器人答案
	GetAnswer(msg)
}

func GetAnswer(msg Msg) {
	token, err1 := getToken(msg)

	if err1 != nil {
		public.Error("Handle getToken error: ", err1.Error())
		return
	}
	//Printf(token)
	msgBytes, err2 := httpPostForm(token)
	if err2 != nil {
		public.Error("\nHandle httpPostForm error: ", err2.Error())
		return
	}
	public.Printf(string(msgBytes))
	var answer XiaoDaiMeng
	err := json.Unmarshal(msgBytes, &answer)
	if err != nil {
		public.Error("\nHandle Unmarshal XiaoDaiMeng error: ", err.Error())
		return
	}

	if answer.AnsNodeId < 1 {
		//回答失败
		public.Printf("\nXiaoDaiMeng 回答失败了->", msg.Content, "<-\n")

		//查找问题列表有没有答案，有接直接用
		hasAnswer := false
		for _, a := range problemList {
			//Printf("strings.Index(msg.Content, a.Problem) ", strings.Index(msg.Content, a.Problem))
			if strings.Index(msg.Content, a.Problem) != -1 {
				answer.Answer = a.Answer
				hasAnswer = true
				break
			}
		}
		if !hasAnswer {
			//随机回复图片
			index := public.GenerateRangeNum(0,8) //
			picName := public.GetCurrentDirectory() + "/static/img/unknow/" + strconv.Itoa(index) + ".jpg"
			public.Debug(picName)
			SendMsg(msg.Sender, picName, PIC_MSG)
			return
		}
	} else {
		if answer.AnswerType == NewsAnswer {
			//新闻
			moreInfo := MoreInfoData{}
			err  = json.Unmarshal([]byte(answer.MoreInfo.NewsAnsDetail),&moreInfo)
			if err != nil {
				public.Error("\nHandle Unmarshal NewsAnsDetail error: ", err.Error())
				return
			}
			news := ""
			for i,m := range moreInfo.MoreInfo.NewsDocs {
				if i < 6 {
					news += m.Cate1 + " " + m.Title + "\n\n"
					news += m.AbsL

					if i<4 {
						news += "\n  -------------  \n"
					}
				}
			}
			public.Debug("\n\n\n\n\n\n\n\n",answer.MoreInfo.NewsAnsDetail)
			if news != "" {
				SendMsg(msg.Sender, news, TXT_MSG)
			} else {
				SendMsg(msg.Sender, FailText, TXT_MSG)
			}
			return
		}
	}

	//Printf("\nCODE :", answer.AnsNodeId)

	//替换机器人的名字
	answer.Answer = strings.Replace(answer.Answer, "小微", XiaoDaiMengName, -1)
	//Printf("\n", answer.Answer)

	go func(a XiaoDaiMeng) {
		SendMsg(a.FromUserName, a.Answer, TXT_MSG)
	}(answer)
}

func SendMsg(wxId string, content string, mType int) {
	lock.Lock()
	var msg = Msg{
		Id:      getTimeId(),
		WxId:    wxId,
		Content: content,
		Type:    mType,
	}
	bRMsg, err := json.Marshal(msg)
	if err != nil {
		public.Error("\nSendMsg Marshal RMsg error: ", err.Error())
		return
	}
	public.Debug(string(bRMsg))
	RWsMsg = make(chan []byte)
	RWsMsg <- bRMsg
	defer lock.Unlock()
}
