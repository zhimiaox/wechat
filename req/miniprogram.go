/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package req

import (
	"fmt"
)

// ListMiniAppParam 列表
type ListMiniAppParam struct {
	PageParam
	AppID      string // 小程序appid
	PlatformID string // 开放平台ID
	Version    string // 版本号
	State      int8   // -1-授权失效 1授权成功，2审核中，3审核通过，4审核失败，5已发布
	OrderBase
}

// MiniAppConfigParam 小程序配置
type MiniAppConfigParam struct {
	AppID       string `binding:"required"` // 小程序appid
	MchID       string // 支付商户号id
	Secret      string // 小程序secret
	ExtConfig   string // 小程序扩展配置，发布时会注入至ext.json
	AutoAudit   int8   // 自动提审(升级) -1否 1是
	AutoRelease int8   // 自动发布-1否 1是
}

type WxACodeParam struct {
	MiniProgramID string `binding:"required"` // 小程序appid
	Page          string
	Scene         string `binding:"omitempty,min=1,max=32"`
}

// Order 获取排序
func (p *ListMiniAppParam) Order() string {
	switch p.OrderColumn {
	case "Version":
		p.OrderColumn = "version"
	default:
		p.OrderColumn = ""
	}
	if p.OrderColumn != "" {
		return fmt.Sprintf("`%s` %s", p.OrderColumn, p.GetSafeClase())
	}
	return "create_time desc"
}
