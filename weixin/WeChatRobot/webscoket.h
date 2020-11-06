#pragma once
#include <string>
//typedef void (*CallbackFun)(const std::string& msgData);
//int initWs(CallbackFun callable);
void startWSs();
void wsSend(const std::string  msg);