/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package api

import (
	"gitee.com/zhimiao/wechat/common/utils"
	"gitee.com/zhimiao/wechat/req"
	"gitee.com/zhimiao/wechat/resp"
	ws "gitee.com/zhimiao/wechat/wechat_service"
	"github.com/gin-gonic/gin"
	"github.com/siddontang/go/log"
)

type payApi struct{}

var PayApi = &payApi{}

func (m *payApi) Test(c *gin.Context) {
	var id string
	if err := c.ShouldBind(&id); err != nil {
		resp.NewApiResult(-4, utils.Validator(err)).Json(c)
		return
	}
	resp.NewApiResult(1).Json(c)
	return
}

// @Summary 异步回调
// @Param PayMchID path string true "微信支付商户号"
// @Router /pay/{PayMchID}/notify [post]
func (m *payApi) Notify(c *gin.Context) {
	config := ws.GetWechatPayConfig(c.GetString(req.PayMchID))
	wechatService := ws.GetWechatService(config)
	if wechatService == nil {
		return
	}
	// 传入request和responseWriter
	server := wechatService.GetServer(c.Request, c.Writer)
	server.SetDebug(true)
	//设置接收消息的处理方法
	server.SetPayHandler(ws.PayHandler)
	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		log.Warn(err)
		return
	}
	//发送回复的消息
	server.Send()
}
