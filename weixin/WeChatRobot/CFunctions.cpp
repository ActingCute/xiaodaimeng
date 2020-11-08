// CFunctions.cpp: 实现文件
//

#include "stdafx.h"
#include "WeChatRobot.h"
#include "CFunctions.h"
#include "afxdialogex.h"
#include "CInformation.h"
#include "CDecryptImage.h"
#include "CMultiOpen.h"
#include "CAddUser.h"
#include "CFriendList.h"
#include "COpenUrl.h"

extern BOOL isAttentTuLing;
BOOL bAutoChat = FALSE;

// CFunctions 对话框

IMPLEMENT_DYNAMIC(CFunctions, CDialogEx)

CFunctions::CFunctions(CWnd* pParent /*=nullptr*/)
	: CDialogEx(IDD_FUNCTIONS, pParent)
{

}

CFunctions::~CFunctions()
{
}

void CFunctions::DoDataExchange(CDataExchange* pDX)
{
	CDialogEx::DoDataExchange(pDX);
}


BEGIN_MESSAGE_MAP(CFunctions, CDialogEx)
	ON_BN_CLICKED(IDC_DECRYPT_PIC, &CFunctions::OnBnClickedDecryptPic)
END_MESSAGE_MAP()




//************************************************************
// 函数名称: OnBnClickedDecryptPic
// 函数说明: 响应解密图片按钮
// 作    者: GuiShou
// 时    间: 2019/7/6
// 参    数: void
// 返 回 值: void
//***********************************************************
void CFunctions::OnBnClickedDecryptPic()
{
	CDecryptImage decryptimage;
	decryptimage.DoModal();
}


