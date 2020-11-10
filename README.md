# 小呆萌

## 一个智障机器人

## go 代理
 
   set GOPROXY=https://goproxy.cn
   
## 配置文件

   go run main.go -c="data/config.json"
   
## 微信
项目来自[WeChatRobot](https://github.com/TonyChen56/WeChatRobot.git)
这个项目是c++的，我只做了暴露一个webscoket端口
编译的exe在 weixin/Release

## 发送消息

    {"m_wxid":"wxid_u3q162gfuq8k22","m_Content":"小呆萌开机了"}
    
## 接受消息

    {"times":"2020-11-10 10-52-06","type":"文字","source":"好友消息","wxid":"wxid_u3q162gfuq8k22","msgSender":"","content":"爱我不"}
    