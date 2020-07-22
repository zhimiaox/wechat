/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package wechat_service

import (
	"gitee.com/zhimiao/wechat-sdk/message"
	"gitee.com/zhimiao/wechat-sdk/pay"
	"gitee.com/zhimiao/wechat/models"
	"github.com/sirupsen/logrus"
	"time"
)

type msgHandler message.MixMessage
type payHandler pay.NotifyResult

// MessageHandler 消息钩子
func MessageHandler(msg message.MixMessage) *message.Reply {
	reply := &message.Reply{
		ReplyScene:   message.ReplySceneOpen,
		ResponseType: message.ResponseTypeString,
	}
	h := msgHandler(msg)
	// 第三方平台回调
	switch h.InfoType {
	// ticket
	case message.InfoTypeVerifyTicket:
		h.TicketHandler(reply)
	// 授权
	case message.InfoTypeAuthorized:
		h.OpenAuthorized(reply)
	// 更新授权
	case message.InfoTypeUpdateAuthorized:
		h.OpenAuthorized(reply)
	// 取消授权
	case message.InfoTypeUnauthorized:
		h.OpenUnAuthorized(reply)
	}
	// 事件消息回调
	if h.MsgType == message.MsgTypeEvent {
		switch h.Event {
		case message.EventWeappAuditSuccess: // 审核通过
			h.MiniAppAuditSuccessHandler(reply)
		case message.EventWeappAuditFail: // 审核失败
			h.MiniAppAuditFailHandler(reply)
		case message.EventWeappAuditDelay: // 审核延后
			h.MiniAppAuditDelayHandler(reply)
		}
	}
	return reply
}

// PayHandler 支付钩子
func PayHandler(msg pay.NotifyResult) *message.Reply {
	reply := &message.Reply{
		ReplyScene:   message.ReplyScenePay,
		ResponseType: message.ResponseTypeXML,
		MsgData: pay.NotifyResp{
			ReturnCode: "SUCCESS",
			ReturnMsg:  "OK",
		},
	}
	h := payHandler(msg)
	// TODO: 完成回调操作，并二次通知具体业务模块
	switch h.PayNotifyInfo {
	case pay.PayTypePay:
	case pay.PayTypeRefund:

	}
	return reply
}

// TicketHandler 票据回调
func (h *msgHandler) TicketHandler(reply *message.Reply) {
	// 记录回调错误
	go func(id string) {
		time.Sleep(5 * time.Second)
		openService := GetOpen(id)
		vo := GetOPState(h.AppID)
		_, err := openService.GetComponentVerifyTicket()
		vo.Ticket = ""
		if err != nil {
			vo.Ticket = err.Error()
		}
		RefreshOPState(h.AppID, vo)
	}(h.AppID)
	return
}

// OpenAuthorized 开放平台授权
func (h *msgHandler) OpenAuthorized(reply *message.Reply) {
	openService := GetOpen(h.AppID)
	info, err := openService.QueryAuthCode(h.AuthorizationCode)
	if err != nil {
		logrus.Error("获取小程序授权信息失败, %s", err.Error())
		return
	}
	authorizerInfo, authorizationInfo, err := openService.GetAuthrInfo(info.Appid)
	if err != nil {
		logrus.Error("获取小程序授权详细信息失败, %s", err.Error())
		return
	}
	miniProgram := &models.Miniprogram{
		AppID:        info.Appid,
		PlatformID:   h.AppID,
		RefreshToken: authorizationInfo.RefreshToken,
		OriginalID:   authorizerInfo.UserName,
	}
	if !miniProgram.Save() {
		logrus.Error("小程序授权失败\n %#v", miniProgram)
	}
}

// OpenUnAuthorized 取消授权
func (h *msgHandler) OpenUnAuthorized(reply *message.Reply) {
	miniProgram := &models.Miniprogram{
		AppID:      h.AuthorizerAppid,
		PlatformID: h.AppID,
		State:      -1,
	}
	if !miniProgram.UpdateByAppID() {
		logrus.Warn("小程序取消授权记录失败\n %#v", miniProgram)
	}
}

// MiniAppAuditSuccessHandler 小程序审核通过
func (h *msgHandler) MiniAppAuditSuccessHandler(reply *message.Reply) {
	var originId string = string(h.ToUserName)
	audit := &models.MiniprogramAudit{
		OriginalID: originId,
		State:      2,
		Reason:     "",
		ScreenShot: "",
	}
	err := audit.ChangeState()
	if err != nil {
		logrus.Warn("小程序[%s]回调状态修改异常, %s", originId, err.Error())
		return
	}
	miniApp := models.Miniprogram{
		OriginalID: originId,
	}
	err = miniApp.GetBySelectKey()
	if err != nil {
		logrus.Warn("小程序[%s]数据获取失败, %s", originId, err.Error())
		return
	}
	// 是否自动发布
	if miniApp.AutoRelease == 1 {
		err := Release(miniApp.PlatformID, miniApp.AppID)
		if err != nil {
			logrus.Warn("小程序[%s]自动发布失败, %s", originId, err.Error())
		}
	}
}

func (h *msgHandler) MiniAppAuditFailHandler(reply *message.Reply) {
	var originId string = string(h.ToUserName)
	audit := &models.MiniprogramAudit{
		OriginalID: originId,
		State:      3,
		Reason:     h.Reason,
		ScreenShot: h.ScreenShot,
	}
	err := audit.ChangeState()
	if err != nil {
		logrus.Warn("小程序[%s]回调状态修改异常, %s", originId, err.Error())
	}
}

// MiniAppAuditDelayHandler 审核延后
func (h *msgHandler) MiniAppAuditDelayHandler(reply *message.Reply) {
	var originId string = string(h.ToUserName)
	audit := &models.MiniprogramAudit{
		OriginalID: originId,
		Reason:     h.Reason,
		ScreenShot: h.ScreenShot,
	}
	err := audit.ChangeState()
	if err != nil {
		logrus.Warn("小程序[%s]回调状态修改异常, %s", originId, err.Error())
	}
}
