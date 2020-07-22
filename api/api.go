/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package api

import (
	"fmt"
	"gitee.com/zhimiao/wechat/common"
	"gitee.com/zhimiao/wechat/common/utils"
	_ "gitee.com/zhimiao/wechat/docs"
	"gitee.com/zhimiao/wechat/req"
	"gitee.com/zhimiao/wechat/resp"
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
	"strconv"
	"time"
)

var Router *gin.Engine

func Start() {
	gin.SetMode(gin.DebugMode)
	// 初始化route
	initRoute()
	httpServer := &http.Server{
		Addr:           common.Config.Server.APIListen,
		Handler:        Router,
		ReadTimeout:    time.Duration(common.Config.Server.ReadTimeOut) * time.Second,
		WriteTimeout:   time.Duration(common.Config.Server.WriteTimeOut) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logrus.Infof("Start HTTP Service Listening %s", common.Config.Server.APIListen)
	httpServer.ListenAndServe()
}

func initRoute() {
	Router = gin.New()
	Router.Use(gin.Recovery(), logMiddleware())
	// 状态监控
	Router.Use(ginprom.PromMiddleware(nil))
	Router.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))
	// 跨域支持
	Router.Use(corsMiddleware())

	/* ------ 文档模块 ------- */
	Router.GET("/docs/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "NAME_OF_ENV_VARIABLE"))

	/* ------ 微信支付模块 ------- */
	pay := Router.Group(fmt.Sprintf("/pay/:%s", req.PayMchID))
	{
		// 注入商户ID
		pay.Use(payMchIDMiddleware())
		pay.Any("/notify", PayApi.Notify)
		pay.POST("/Test", PayApi.Test)
	}

	/* ------ 小程序模块 ------- */
	miniApp := Router.Group(fmt.Sprintf("/miniprogram"))
	{
		miniApp.GET("/Lists", Miniprogram.Lists)
		miniApp.POST("/Config", Miniprogram.Config)
		miniApp.GET("/GetWXACodeUnlimit", Miniprogram.GetWXACodeUnlimit)
		miniApp.GET("/Send", Miniprogram.Send)
	}

	/* ------ 微信开放平台模块 ------- */
	oplatform := Router.Group(fmt.Sprintf("/oplatform/:%s", req.PlatformID))
	{
		// 默认manage
		oplatform.GET("/Lists", OPlatform.Lists)
		oplatform.POST("/Set", OPlatform.Set)

		// 注入平台ID
		oplatform.Use(oplatformIDMiddleware())

		oplatform.Any("/notify", OPlatform.Notify)
		oplatform.GET("/auth", OPlatform.Auth)
		oplatform.GET("/redirect", OPlatform.Redirect)

		oplatform.GET("/tpl/draft", OPlatform.TplDraft)
		oplatform.GET("/tpl/list", OPlatform.TplList)
		oplatform.POST("/tpl/pushToAuto", OPlatform.pushToAuto)
		oplatform.POST("/tpl/add", OPlatform.TplAdd)
		oplatform.DELETE("/tpl/del", OPlatform.TplDel)

		oplatform.GET("/GetAuthorizerInfo", OPlatform.GetAuthorizerInfo)
		oplatform.GET("/GetTestQrcode", OPlatform.GetTestQrcode)
		oplatform.GET("/GetCodeCategory", OPlatform.GetCodeCategory)
		oplatform.GET("/getCodePageList", OPlatform.GetCodePageList)

		oplatform.GET("/GetLatestAuditStatus", OPlatform.GetLatestAuditStatus)
		oplatform.POST("/SetDomain", OPlatform.SetDomain)
		oplatform.POST("/CommitCode", OPlatform.CommitCode)
		oplatform.POST("/SubmitAudit", OPlatform.SubmitAudit)
		oplatform.POST("/UndoCodeAudit", OPlatform.UndoCodeAudit)
		oplatform.POST("/SpeedUpAudit", OPlatform.SpeedUpAudit)
		oplatform.POST("/Release", OPlatform.Release)
	}

}

// corsMiddleware 跨域
func corsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"*"},
	})
}

// logMiddleware 日志中间件
func logMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUrl := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求ip
		clientIP := c.ClientIP()
		// 日志格式
		logrus.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUrl,
		}).Info()
	}
}

// jwtMiddleware jwt鉴权
func jwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var msg string = "登陆信息获取失败"
		token := c.GetHeader("LOGIN-KEY")
		if token != "" {
			str, err := utils.ParseToken(token)
			if err != nil {
				print(err.Error())
				msg = "登录验证失败"
			}
			if str != "" {
				c.Set("LOGIN-TOKEN", str)
				uid, _ := strconv.Atoi(str)
				c.Set("UID", uid)
				if uid > 0 {
					msg = ""
				}
			}
		}
		if msg != "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"msg":  msg,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 第三方平台模块中间件 - 平台ID注入
func oplatformIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param(req.PlatformID)
		if id == "" {
			id = c.Query(req.PlatformID)
		}
		if id == "" {
			resp.NewApiResult(-4, "无法识别").Json(c)
			c.Abort()
			return
		}
		c.Set(req.PlatformID, id)
		c.Next()
	}
}

// 支付模块商户号注入
func payMchIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param(req.PayMchID)
		if id == "" {
			id = c.Query(req.PayMchID)
		}
		if id == "" {
			resp.NewApiResult(-4, "无法识别").Json(c)
			c.Abort()
			return
		}
		c.Set(req.PayMchID, id)
		c.Next()
	}
}
