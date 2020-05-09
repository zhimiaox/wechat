/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

// MiniprogramAudit 微信小程序提交审核的小程序
type MiniprogramAudit struct {
	ID                   int       `gorm:"primary_key;column:id;type:int(10) unsigned;not null"`     // ID
	AppID                string    `gorm:"column:app_id;type:varchar(64);not null"`                  // 小程序appId
	OriginalID           string    `gorm:"column:original_id;type:varchar(45);not null"`             // 小程序原始id
	AuditID              uint64    `gorm:"column:audit_id;type:bigint(20) unsigned;not null"`        // 审核编号
	State                int8      `gorm:"column:state;type:tinyint(3);not null"`                    // 审核状态，-1-撤销审核 1为审核中，2为审核成功，3为审核失败
	Reason               string    `gorm:"column:reason;type:varchar(255)"`                          // 当status=1，审核被拒绝时，返回的拒绝原因
	ScreenShot           string    `gorm:"column:screen_shot;type:varchar(3000)"`                    // 附件素材
	TemplateID           int       `gorm:"column:template_id;type:int(11);not null"`                 // 最新提交审核或者发布的模板id
	TemplateAppID        string    `gorm:"column:template_app_id;type:varchar(64);not null"`         // 模板开发小程序ID
	TemplateAppName      string    `gorm:"column:template_app_name;type:varchar(255)"`               // 开发小程序名
	TemplateAppDeveloper string    `gorm:"column:template_app_developer;type:varchar(255);not null"` // 开发者
	TemplateDesc         string    `gorm:"column:template_desc;type:varchar(64);not null"`           // 模板描述
	TemplateVersion      string    `gorm:"column:template_version;type:varchar(64);not null"`        // 模板版本号
	CreateTime           time.Time `gorm:"column:create_time;type:datetime"`
	UpdateTime           time.Time `gorm:"column:update_time;type:datetime"`
}

// GetBySelectKey 根据ID获取平台数据
func (m *MiniprogramAudit) GetBySelectKey() (err error) {
	db := Mysql.Model(m)
	if m.AppID != "" {
		db = db.Where("app_id=?", m.AppID)
	}
	if m.AuditID != 0 {
		db = db.Where("audit_id=?", m.AuditID)
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

// Change 审核变更
func (m *MiniprogramAudit) ChangeState() (err error) {
	var db *gorm.DB
	// 审核中-创建
	if m.State == 1 {
		db = Mysql.Model(m).Create(m)
		if db.Error != nil || db.RowsAffected == 0 {
			err = fmt.Errorf("审核单创建失败")
			return
		}
		miniApp := Miniprogram{
			AppID:         m.AppID,
			State:         2,
			AuditID:       m.AuditID,
			Version:       m.TemplateVersion,
			NowTemplateID: m.TemplateID,
		}
		db = Mysql.Model(miniApp).Update(miniApp)
		if db.Error != nil || db.RowsAffected == 0 {
			err = fmt.Errorf("小程序状态更新失败")
			return
		}
	} else {
		p := *m
		m.GetBySelectKey()
		p.ID = m.ID
		if m.State != 1 {
			err = fmt.Errorf("检索最新审核单状态异常")
			return
		}
		db = Mysql.Model(m).Update(&p)
		if db.Error != nil || db.RowsAffected == 0 {
			err = fmt.Errorf("审核单更新失败")
			return
		}
		miniApp := Miniprogram{
			AppID: m.AppID,
		}
		switch m.State {
		case -1: // 撤审
			miniApp.State = 6
		case 2: // 成功
			miniApp.State = 3
		case 3: // 失败
			miniApp.State = 4
			miniApp.AutoAudit = -1
		}
		if miniApp.State != 0 && !miniApp.UpdateByAppID() {
			err = fmt.Errorf("小程序状态更新失败")
		}
	}
	return
}
