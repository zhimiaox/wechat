/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package resp

import "time"

// PlatformStateVO 平台状态检测
type PlatformStateVO struct {
	Ticket       string // 令牌状态
	Notify       string // 回调状态
	VDomainCheck string // 域名状态
}

// Platform 平台注册信息主表
type PlatformListsVO struct {
	PlatformID      string // 平台 appid
	PlatformSecret  string // 平台 appsecret
	PlatformToken   string // 平台 token
	PlatformKey     string // 平台 消息解密key
	ServerDomain    string // 服务器域名
	BizDomain       string // 业务域名
	AuthRedirectURL string // 用户授权成功回跳地址
	CreateTime      time.Time
	UpdateTime      time.Time
	VPlatformState  PlatformStateVO
}
