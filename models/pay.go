/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package models

import (
	"time"
)

// Pay 微信支付
type Pay struct {
	MchID        string    `gorm:"primary_key;column:mch_id;type:varchar(32);not null"` // 支付商户号id
	Token        string    `gorm:"column:token;type:varchar(255)"`                      // 支付密钥
	Cert         string    `gorm:"column:cert;type:longtext"`                           // 支付证书
	PayNotifyURL string    `gorm:"column:pay_notify_url;type:varchar(255)"`             // 支付回调
	PayRefundURL string    `gorm:"column:pay_refund_url;type:varchar(255)"`             // 退款回调
	CreateTime   time.Time `gorm:"column:create_time;type:datetime"`
	UpdateTime   time.Time `gorm:"column:update_time;type:datetime"`
}

func (m *Pay) GetByMchId() (rows int64, err error) {
	db := Mysql.Where("mch_id=?", m.MchID).First(m)
	err = db.Error
	rows = db.RowsAffected
	return
}
