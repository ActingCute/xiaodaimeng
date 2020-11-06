// WebSocketTest.cpp : �������̨Ӧ�ó������ڵ㡣
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

//������ֵ���ݶ�
struct NameAndValue
{
	std::string strName;
	std::string strValue;
};
// �ַ����ָ�
int StringSplit(std::vector<std::string>& dst, const std::string& src, const std::string& separator);
//ȥǰ��ո�
std::string& StringTrim(std::string& str);
//��ȡ�������������
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
	//����websocket upgrade�ɹ�֮�󣬵���open_handler�������ص�on_open��
	//��������Ի�ȡhttp����ĵ�ַ��������Ϣ��
	std::cout << "open handler" << std::endl;
	/*
	server::connection_ptr con = s->get_con_from_hdl(hdl);
	websocketpp::config::core::request_type requestClient = con->get_request();
	std::string strMethod = requestClient.get_method();		//���󷽷�
	std::string strUri = requestClient.get_uri();			//����uri��ַ�����Խ�������
	std::string strRequestOperateCommand = "";				//��������
	std::list<NameAndValue> listRequestOperateParameter;	//���������б�
	GetReqeustCommandAndParmeter(strUri, strRequestOperateCommand, listRequestOperateParameter);
	std::cout << "command:" << strRequestOperateCommand << std::endl;
	*/
	//s->send(hdl,"�ٺ����ֳɹ�", websocketpp::frame::opcode::TEXT);

}

void on_close(websocketpp::connection_hdl hdl) {
	std::cout << "Close handler" << std::endl;
}
using namespace rapidjson;		//ʹ��rapidjson�����ռ䣬���û�������䣬���޷�ʹ��Document���������ͣ�

std::string asyncSendWsToWx(std::string& msg) {
	try {

		char* json = (char*)msg.c_str();
		Document document;
		document.Parse(json);

		//���ַ����ж�ȡ����  
		if (document["m_wxid"].GetString() && document["m_Content"].GetString())
		{
			std::string m_wxid = document["m_wxid"].GetString();
			std::string m_Content = document["m_Content"].GetString();

			wchar_t* w_m_wxid = StringToWchar_t(m_wxid);
			wchar_t* w_m_Content = StringToWchar_t(m_Content);


			//������ݵ��ṹ��
			MessageStruct* message = new MessageStruct;
			wcscpy_s(message->wxid, wcslen(w_m_wxid) + 1, w_m_wxid);
			wcscpy_s(message->content, wcslen(w_m_Content) + 1, w_m_Content);

			//���͵�΢��
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
		hdl.lock().get() ������ӱ�ʶ
		msg->get_payload() ���յ�����Ϣ����
		msg->get_opcode() ���յ���Ϣ������ ���������ı�TEXT,������BINARY�ȵ�
	*/
	//std::cout << "on_message called with hdl: " << hdl.lock().get()
	//	<< " and message: " << msg->get_payload()
	//	<< std::endl;

	try {
		/*
			������Ϣ
			s->send(
			hdl, //����
			msg->get_payload(), //��Ϣ
			msg->get_opcode());//��Ϣ����
		*/
		//s->send(hdl, msg->get_payload(), msg->get_opcode());

		//�ַ���  
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
		std::cout << "�˿� 9100" << std::endl;

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


// �ַ����ָ�
int StringSplit(std::vector<std::string>& dst, const std::string& src, const std::string& separator)
{
	if (src.empty() || separator.empty())
		return 0;

	int nCount = 0;
	std::string temp;
	size_t pos = 0, offset = 0;

	// �ָ��1~n-1��
	while ((pos = src.find_first_of(separator, offset)) != std::string::npos)
	{
		temp = src.substr(offset, pos - offset);
		if (temp.length() > 0) {
			dst.push_back(temp);
			nCount++;
		}
		offset = pos + 1;
	}

	// �ָ��n��
	temp = src.substr(offset, src.length() - offset);
	if (temp.length() > 0) {
		dst.push_back(temp);
		nCount++;
	}

	return nCount;
}
//ȥǰ��ո�
std::string& StringTrim(std::string& str)
{
	if (str.empty()) {
		return str;
	}
	str.erase(0, str.find_first_not_of(" "));
	str.erase(str.find_last_not_of(" ") + 1);
	return str;
}
//��ȡ�������������
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