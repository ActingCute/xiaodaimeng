// WebSocketTest.cpp : 定义控制台应用程序的入口点。
//

#include "stdafx.h"

#include <boost/algorithm/string.hpp>
#include <string>
#include <vector>
#include <list>

#include <iostream>
#include <websocketpp/config/asio_no_tls.hpp>
#include <websocketpp/server.hpp>
#include "webscoket.h"

#include <future>
#include "CSendMsg.h"
#include "CChatRecords.h"

//json
#include "rapidjson/document.h"
#include <vector>
#include "rapidjson/stringbuffer.h"
#include "rapidjson/writer.h"

//名称与值数据对
struct NameAndValue
{
	std::string strName;
	std::string strValue;
};
// 字符串分割
int StringSplit(std::vector<std::string>& dst, const std::string& src, const std::string& separator);
//去前后空格
std::string& StringTrim(std::string& str);
//获取请求命令与参数
bool GetReqeustCommandAndParmeter(std::string strUri, std::string& strRequestOperateCommand, std::list<NameAndValue>& listRequestOperateParameter);




typedef websocketpp::server<websocketpp::config::asio> server;

using websocketpp::lib::placeholders::_1;
using websocketpp::lib::placeholders::_2;
using websocketpp::lib::bind;

// pull out the type of messages sent by our config
typedef server::message_ptr message_ptr;

bool validate(server*, websocketpp::connection_hdl) {
	//sleep(6);
	return true;
}

void on_http(server* s, websocketpp::connection_hdl hdl) {
	server::connection_ptr con = s->get_con_from_hdl(hdl);

	std::string res = con->get_request_body();

	std::stringstream ss;
	ss << "got HTTP request with " << res.size() << " bytes of body data.";

	con->set_body(ss.str());
	con->set_status(websocketpp::http::status_code::ok);
}

void on_fail(server* s, websocketpp::connection_hdl hdl) {
	server::connection_ptr con = s->get_con_from_hdl(hdl);

	std::cout << "Fail handler: " << con->get_ec() << " " << con->get_ec().message() << std::endl;
}

server* WsServer;
websocketpp::connection_hdl WsHdl;



void wsSend(const std::string msg) {
	WsServer->send(WsHdl, msg, websocketpp::frame::opcode::text);
}

void on_open(server* s, websocketpp::connection_hdl hdl) {
	WsServer = s;
	WsHdl = hdl;
	//申请websocket upgrade成功之后，调用open_handler函数，回调on_open。
	//在这里，可以获取http请求的地址、参数信息。
	std::cout << "open handler" << std::endl;
	/*
	server::connection_ptr con = s->get_con_from_hdl(hdl);
	websocketpp::config::core::request_type requestClient = con->get_request();
	std::string strMethod = requestClient.get_method();		//请求方法
	std::string strUri = requestClient.get_uri();			//请求uri地址，可以解析参数
	std::string strRequestOperateCommand = "";				//操作类型
	std::list<NameAndValue> listRequestOperateParameter;	//操作参数列表
	GetReqeustCommandAndParmeter(strUri, strRequestOperateCommand, listRequestOperateParameter);
	std::cout << "command:" << strRequestOperateCommand << std::endl;
	*/
	//s->send(hdl,"嘿嘿握手成功", websocketpp::frame::opcode::TEXT);

}

void on_close(websocketpp::connection_hdl hdl) {
	std::cout << "Close handler" << std::endl;
}
using namespace rapidjson;		//使用rapidjson命名空间，如果没有这个语句，将无法使用Document等数据类型，

std::string asyncSendWsToWx(std::string& msg) {
	try {

		char* json = (char*)msg.c_str();
		Document document;
		document.Parse(json);

		//从字符串中读取数据  
		if (document["m_wxid"].GetString() && document["m_Content"].GetString())
		{
			std::string m_wxid = document["m_wxid"].GetString();
			std::string m_Content = document["m_Content"].GetString();

			wchar_t* w_m_wxid = StringToWchar_t(m_wxid);
			wchar_t* w_m_Content = StringToWchar_t(m_Content);


			//填充数据到结构体
			MessageStruct* message = new MessageStruct;
			wcscpy_s(message->wxid, wcslen(w_m_wxid) + 1, w_m_wxid);
			wcscpy_s(message->content, wcslen(w_m_Content) + 1, w_m_Content);

			//发送到微信
			webSendMsg(message);
		}


	}
	catch (const std::exception& e) {

	}

	return "";
}




// Define a callback to handle incoming messages
void on_message(server* s, websocketpp::connection_hdl hdl, message_ptr msg) {
	/*
		hdl.lock().get() 获得连接标识
		msg->get_payload() 是收到的消息内容
		msg->get_opcode() 是收到消息的类型 ，包含：文本TEXT,二进制BINARY等等
	*/
	//std::cout << "on_message called with hdl: " << hdl.lock().get()
	//	<< " and message: " << msg->get_payload()
	//	<< std::endl;

	try {
		/*
			发送消息
			s->send(
			hdl, //连接
			msg->get_payload(), //消息
			msg->get_opcode());//消息类型
		*/
		//s->send(hdl, msg->get_payload(), msg->get_opcode());

		//字符串  
		auto futureFunction = std::async(asyncSendWsToWx, msg->get_payload());

	}
	catch (websocketpp::exception const& e) {
		std::cout << "Echo failed because: "
			<< "(" << e.what() << ")" << std::endl;
	}
}


void startWSs() {
	server print_server;
	try {
		// Set logging settings
		print_server.set_access_channels(websocketpp::log::alevel::all);
		print_server.set_error_channels(websocketpp::log::elevel::all);
		//print_server.clear_access_channels(websocketpp::log::alevel::frame_payload);

		// Register our message handler
		print_server.set_message_handler(bind(&on_message, &print_server, ::_1, ::_2));
		print_server.set_http_handler(bind(&on_http, &print_server, ::_1));
		print_server.set_fail_handler(bind(&on_fail, &print_server, ::_1));
		print_server.set_open_handler(bind(&on_open, &print_server, ::_1));
		print_server.set_close_handler(bind(&on_close, ::_1));
		print_server.set_validate_handler(bind(&validate, &print_server, ::_1));

		// Initialize ASIO
		print_server.init_asio();
		print_server.set_reuse_addr(true);

		// Listen on port 9100
		print_server.listen(9100);
		// Start the server accept loop
		print_server.start_accept();
		std::cout << "端口 9100" << std::endl;

		// Start the ASIO io_service run loop
		print_server.run();

		//stop
		//print_server.stop();
	}
	catch (websocketpp::exception const& e) {
		std::cout << e.what() << std::endl;
	}
	catch (const std::exception& e) {
		std::cout << e.what() << std::endl;
	}
	catch (...) {
		std::cout << "other exception" << std::endl;
	}
}


// 字符串分割
int StringSplit(std::vector<std::string>& dst, const std::string& src, const std::string& separator)
{
	if (src.empty() || separator.empty())
		return 0;

	int nCount = 0;
	std::string temp;
	size_t pos = 0, offset = 0;

	// 分割第1~n-1个
	while ((pos = src.find_first_of(separator, offset)) != std::string::npos)
	{
		temp = src.substr(offset, pos - offset);
		if (temp.length() > 0) {
			dst.push_back(temp);
			nCount++;
		}
		offset = pos + 1;
	}

	// 分割第n个
	temp = src.substr(offset, src.length() - offset);
	if (temp.length() > 0) {
		dst.push_back(temp);
		nCount++;
	}

	return nCount;
}
//去前后空格
std::string& StringTrim(std::string& str)
{
	if (str.empty()) {
		return str;
	}
	str.erase(0, str.find_first_not_of(" "));
	str.erase(str.find_last_not_of(" ") + 1);
	return str;
}
//获取请求命令与参数
bool GetReqeustCommandAndParmeter(std::string strUri, std::string& strRequestOperateCommand, std::list<NameAndValue>& listRequestOperateParameter)
{
	bool bRet = false;
	std::vector<std::string> vecRequest;
	int nRetSplit = StringSplit(vecRequest, strUri, "?");
	if (nRetSplit > 0)
	{
		if (vecRequest.size() == 1)
		{
			strRequestOperateCommand = vecRequest[0];
		}
		else if (vecRequest.size() > 1)
		{
			strRequestOperateCommand = vecRequest[0];
			std::string strRequestParameter = vecRequest[1];
			std::vector<std::string> vecParams;
			nRetSplit = StringSplit(vecParams, strRequestParameter, "&");
			if (nRetSplit > 0)
			{
				std::vector<std::string>::iterator iter, iterEnd;
				iter = vecParams.begin();
				iterEnd = vecParams.end();
				for (iter; iter != iterEnd; iter++)
				{
					std::vector<std::string> vecNameOrValue;
					nRetSplit = StringSplit(vecNameOrValue, *iter, "=");
					if (nRetSplit > 0)
					{
						NameAndValue nvNameAndValue;
						nvNameAndValue.strName = vecNameOrValue[0];
						nvNameAndValue.strValue = "";
						if (vecNameOrValue.size() > 1)
						{
							nvNameAndValue.strValue = vecNameOrValue[1];
						}
						//insert
						listRequestOperateParameter.push_back(nvNameAndValue);
					}
				}
			}
		}
		else
		{

		}
	}
	return bRet;
}