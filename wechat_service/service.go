/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

/*
	服务模块，用于获取各种操作句柄
*/
package wechat_service

import (
	"gitee.com/zhimiao/wechat-sdk"
	"gitee.com/zhimiao/wechat-sdk/miniprogram"
	"gitee.com/zhimiao/wechat-sdk/open"
	"gitee.com/zhimiao/wechat-sdk/pay"
	"gitee.com/zhimiao/wechat/models"
)

// GetWechatService 初始化获取微信原始句柄
func GetWechatService(cfg *Config) *wechat.Wechat {
	if cfg == nil {
		return nil
	}
	if cfg.Cache == nil {
		cfg.Cache = models.Redis
	}
	a := wechat.Config(*cfg)
	return wechat.NewWechat(&a)
}

// GetOpen 根据appid获取开放平台句柄
func GetOpen(platformID string) *open.Open {
	config := GetOpenConfig(platformID)
	if config == nil {
		return nil
	}
	return GetWechatService(config).GetOpen()
}

// GetOpenMiniPrograms 获取开放平台掌控小程序句柄
func GetOpenMiniPrograms(platformId, MiniProgramID string) (openMaService *open.MiniPrograms) {
	wechatService := GetOpen(platformId)
	if wechatService == nil {
		return
	}
	refreshToken := GetOpenMaRefreshToken(MiniProgramID)
	if refreshToken == "" {
		return
	}
	openMaService = wechatService.NewMiniPrograms(MiniProgramID, refreshToken)
	return
}

// GetWechatPay 根据商户号获取微信支付句柄
func GetWechatPay(mchId string) *pay.Pay {
	config := GetWechatPayConfig(mchId)
	return GetWechatService(config).GetPay()
}

// GetMiniApp 根据appid获取小程序句柄
func GetMiniApp(appID string) *miniprogram.MiniProgram {
	config := GetMiniAppConfig(appID)
	return GetWechatService(config).GetMiniProgram()
}
