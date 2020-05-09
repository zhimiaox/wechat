/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package req

// 平台编辑
type PlatformParam struct {
	PlatformID      string `binding:"required"` // 平台 appid
	PlatformSecret  string // 平台 appsecret
	PlatformToken   string // 平台 token
	PlatformKey     string // 平台 消息解密key
	ServerDomain    string // 服务器域名
	BizDomain       string // 业务域名
	AuthRedirectURL string // 用户授权成功回跳地址
}

// 提交代码
type CommitCodeParam struct {
	MiniProgramID string `binding:"required"` // 小程序ID
	TemplateID    int    `binding:"required"` // 模板编号
	// ExtMap        map[string]string // 扩展
	// UserVersion   string            `binding:"required"` // 提交版本
	// UserDesc      string            `binding:"required"` // 版本说明
}

// SpeedUpAuditParam 加急审核
type SpeedUpAuditParam struct {
	MiniProgramID string `binding:"required"` // 小程序ID
	AuditID       uint64 `binding:"required"` // 要加急的审核单号
}
