/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package wechat_service

import (
	"encoding/base64"
	"fmt"

	"github.com/zhi-miao/wechat-sdk"
	"github.com/zhi-miao/wechat/common"
	"github.com/zhi-miao/wechat/models"
)

type Config wechat.Config

// GetOpenConfig 获取微信开放平台配置
func GetOpenConfig(platformID string) (config *Config) {
	m := &models.Platform{
		PlatformID: platformID,
	}
	rows, err := m.GetByAppId()
	if err != nil || rows == 0 {
		return
	}
	config = &Config{
		AppID:          m.PlatformID,
		AppSecret:      m.PlatformSecret,
		Token:          m.PlatformToken,
		EncodingAESKey: m.PlatformKey,
	}
	return
}

// GetWechatPayConfig 获取微信支付配置
func GetWechatPayConfig(mchId string) (config *Config) {
	m := &models.Pay{
		MchID: mchId,
	}
	rows, err := m.GetByMchId()
	if err != nil || rows == 0 {
		return
	}
	p12Data, err := base64.StdEncoding.DecodeString(m.Cert)
	if err != nil {
		return
	}
	config = &Config{
		PayMchID:     m.MchID,
		PayNotifyURL: fmt.Sprintf("%s/pay/%s/notify", common.Config.Server.ApiHost, m.MchID),
		PayKey:       m.Token,
		P12:          p12Data,
	}
	return
}

// GetMiniAppConfig 获取小程序配置
func GetMiniAppConfig(appId string) (config *Config) {
	m := &models.Miniprogram{
		AppID: appId,
	}
	rows, err := m.GetByAppID()
	if err != nil || rows == 0 {
		return
	}
	config = &Config{
		AppID:     m.AppID,
		AppSecret: m.Secret,
	}
	return
}

// GetOpenMaRefreshToken 获取代小程序api请求refreshToken
func GetOpenMaRefreshToken(appid string) string {
	maModel := &models.Miniprogram{
		AppID: appid,
	}
	if rows, err := maModel.GetByAppID(); err != nil || rows == 0 {
		return ""
	}
	return maModel.RefreshToken
}
