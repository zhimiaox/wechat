/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package models

import (
	"fmt"
	"gitee.com/zhimiao/wechat/req"
	"time"
)

// Miniprogram 微信小程序授权表列表
type Miniprogram struct {
	AppID          string    `gorm:"primary_key;column:app_id;type:varchar(45);not null"` // 小程序appid
	PlatformID     string    `gorm:"column:platform_id;type:varchar(32);not null"`        // 开放平台ID
	MchID          string    `gorm:"column:mch_id;type:varchar(32)"`                      // 支付商户号id
	OriginalID     string    `gorm:"column:original_id;type:varchar(45);not null"`        // 原始ID
	RefreshToken   string    `gorm:"column:refresh_token;type:varchar(255);not null"`     // 接口调用凭据刷新令牌
	Secret         string    `gorm:"column:secret;type:varchar(50)"`                      // 小程序secret
	ExtConfig      string    `gorm:"column:ext_config;type:text"`                         // 小程序扩展配置
	State          int8      `gorm:"column:state;type:tinyint(3);not null"`               // -1-授权失效 1授权成功，2审核中，3审核通过，4审核失败，5已发布 6已撤审
	Version        string    `gorm:"column:version;type:varchar(30)"`                     // 当前版本
	NowTemplateID  int       `gorm:"column:now_template_id;type:int(10) unsigned"`        // 当前模板ID
	TemplateListen string    `gorm:"column:template_listen;type:varchar(64);not null"`    // 模板监听开发小程序(appid)
	AuditID        uint64    `gorm:"column:audit_id;type:bigint(20) unsigned;not null"`   // 审核编号
	AutoAudit      int8      `gorm:"column:auto_audit;type:tinyint(2);not null"`          // 自动提审(升级) -1否 1是
	AutoRelease    int8      `gorm:"column:auto_release;type:tinyint(2);not null"`        // 自动发布-1否 1是
	CreateTime     time.Time `gorm:"column:create_time;type:datetime"`
	UpdateTime     time.Time `gorm:"column:update_time;type:datetime"`
}

// ListAutoAudit 获取自动升级的小程序列表
func (m *Miniprogram) ListAutoAudit() (data []Miniprogram) {
	db := Mysql.Model(m).Select("app_id")
	db = db.Where("platform_id=? and template_listen=? and now_template_id<? and auto_audit=1 and state>0", m.PlatformID, m.TemplateListen, m.NowTemplateID)
	db.Find(&data)
	return
}

// List 获取小程序列表
func (m *Miniprogram) List(param *req.ListMiniAppParam) (data []Miniprogram, rows int) {
	db := Mysql.Model(m)
	if param.AppID != "" {
		db = db.Where("app_id like ?", "%"+param.AppID+"%")
	}
	if param.PlatformID != "" {
		db = db.Where("platform_id like ?", "%"+param.PlatformID+"%")
	}
	if param.Version != "" {
		db = db.Where("version=?", param.Version)
	}
	if param.State != 0 {
		db = db.Where("state=?", param.State)
	}

	db.Offset(param.Offset()).Limit(param.PageSize).Order(param.Order()).Find(&data).Offset(-1).Limit(-1).Count(&rows)
	return data, rows
}

// GetByAppID 根据appid获取小程序信息
func (m *Miniprogram) GetByAppID() (rows int64, err error) {
	db := Mysql.Model(m).First(m)
	err = db.Error
	rows = db.RowsAffected
	return
}

// GetBySelectKey 根据ID获取小程序数据
func (m *Miniprogram) GetBySelectKey() (err error) {
	db := Mysql.Model(m)
	if m.AppID != "" {
		db = db.Where("app_id=?", m.AppID)
	}
	if m.PlatformID != "" {
		db = db.Where("platform_id=?", m.PlatformID)
	}
	if m.OriginalID != "" {
		db = db.Where("original_id=?", m.OriginalID)
	}
	db.Last(m)
	err = db.Error
	if err == nil && db.RowsAffected == 0 {
		err = fmt.Errorf("查询失败")
	}
	return
}

// Save 保存数据
func (m *Miniprogram) Save() bool {
	var count int
	Mysql.Model(m).Where("app_id=?", m.AppID).Count(&count)
	if count > 0 {
		return m.UpdateByAppID()
	} else {
		db := Mysql.Create(m)
		if db.RowsAffected > 0 {
			return true
		}
	}
	return false
}

// UpdateByAppID 根据小程序ID更新小程序配置
func (m *Miniprogram) UpdateByAppID() bool {
	if m.AppID == "" {
		return false
	}
	db := Mysql.Model(m).Update(m)
	return db.Error == nil && db.RowsAffected > 0
}
