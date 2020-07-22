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

// Platform 平台注册信息主表
type Platform struct {
	PlatformID      string    `gorm:"primary_key;column:platform_id;type:varchar(32);not null"` // 平台 appid
	PlatformSecret  string    `gorm:"column:platform_secret;type:varchar(32)"`                  // 平台 appsecret
	PlatformToken   string    `gorm:"column:platform_token;type:varchar(255)"`                  // 平台 token
	PlatformKey     string    `gorm:"column:platform_key;type:varchar(255)"`                    // 平台 消息解密key
	ServerDomain    string    `gorm:"column:server_domain;type:varchar(1000)"`                  // 服务器域名
	BizDomain       string    `gorm:"column:biz_domain;type:varchar(1000)"`                     // 业务域名
	AuthRedirectURL string    `gorm:"column:auth_redirect_url;type:varchar(300)"`               // 用户授权成功回跳地址
	CreateTime      time.Time `gorm:"column:create_time;type:datetime"`
	UpdateTime      time.Time `gorm:"column:update_time;type:datetime"`
}

// Set 保存平台信息
func (m *Platform) Set() (err error) {
	var (
		db   *gorm.DB = Mysql.Model(m)
		rows int
	)
	db.Count(&rows)
	if rows > 0 {
		db = db.Update(m)
	} else {
		db = db.Create(m)
	}
	err = db.Error
	if err == nil && db.RowsAffected == 0 {
		err = fmt.Errorf("保存失败")
	}
	return
}

// List 获取平台列表
func (m *Platform) List(offset int, limit int) ([]Platform, int) {
	data := []Platform{}
	rows := 0
	db := Mysql.Model(m)
	if m.PlatformID != "" {
		db = db.Where("platform_id like ?", "%"+m.PlatformID+"%")
	}
	db.Offset(offset).Limit(limit).Order("create_time desc").Find(&data).Offset(-1).Limit(-1).Count(&rows)
	return data, rows
}

// GetByAppId 根据ID获取平台数据
func (m *Platform) GetByAppId() (rows int64, err error) {
	db := Mysql.Model(m).First(m)
	err = db.Error
	rows = db.RowsAffected
	return
}

// ExportList 数据输出
func (m *Platform) ExportList(lastId int) ([]Platform, int) {
	data := []Platform{}
	db := Mysql
	if lastId > 0 {
		db = db.Where("id > ?", lastId)
	}
	db.Limit(1000).Order("id asc").Find(&data)
	return data, len(data)
}
