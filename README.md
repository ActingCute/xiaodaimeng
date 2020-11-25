# 小呆萌

## 一个智障机器人

## go 代理
 
   set GOPROXY=https://goproxy.cn
  
## 编译32位

    set GOARCH=386   
   
## vendor

    go mod vendor   
   
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
    
## 注意

    第一次运行，微信会自动退出，在任务管理器删掉再重新来一遍
    
微信版本只支持  `2.6.8.52` 其他版本不支持 

    {"content":"在不","id":"20201125081833","receiver":"self","sender":"wxid_u3q162gfuq8k22","srvid":1,"time":"2020-11-25 08:18:33","type":1}
    {"content":"在不","id":"20201125081149","receiver":"22925504714@chatroom","sender":"self","srvid":1,"time":"2020-11-25 08:11:49","type":1}
    {"content":"在不","id":"20201125083532","receiver":"22925504714@chatroom","sender":"wxid_u3q162gfuq8k22","srvid":1,"time":"2020-11-25 08:35:32","type":1}
