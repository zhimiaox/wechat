/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package resp

import "time"

// MiniAppStateVO 小程序状态检测
type MiniAppStateVO struct {
	VConfigState string // 配置状态
	VPayState    string // 支付有效状态
}

// Miniprogram 微信小程序列表
type MiniAppListsVO struct {
	AppID          string // 小程序appid
	PlatformID     string // 开放平台ID
	MchID          string // 支付商户号id
	OriginalID     string // 原始ID
	RefreshToken   string // 接口调用凭据刷新令牌
	Secret         string // 小程序secret
	ExtConfig      string // 小程序扩展配置，发布时会注入至ext.json
	State          int8   // -1-授权失效 1授权成功，2审核中，3审核通过，4审核失败，5已发布
	Version        string // 当前版本
	TemplateListen string // 模板监听开发小程序(appid)
	AuditID        uint64 // 审核编号
	AutoAudit      int8   // 自动提审(升级) -1否 1是
	AutoRelease    int8   // 自动发布-1否 1是
	CreateTime     time.Time
	UpdateTime     time.Time
}
