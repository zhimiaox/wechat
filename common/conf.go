/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package common

import (
	"bytes"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/configor"
)

type config struct {
	App struct {
		PageSize         int
		JwtSecret        string
		MaxAuditProgress int
	}
	Server struct {
		ApiHost      string
		APIListen    string
		UDPListen    string
		ReadTimeOut  int
		WriteTimeOut int
	}
	Mysql struct {
		Host        string
		User        string
		Password    string
		Database    string
		TablePrefix string
	}
	InfluxDB struct {
		Host  string
		Token string
		Org   string
	}
	Redis struct {
		Host        string
		Auth        string
		MaxIdle     int
		MaxActive   int
		IdleTimeOut int
	}
	Aliyun struct {
		AccessKey    string
		AccessSecret string
		RegionId     string
		SmsConfig    struct {
			SignName     string
			TemplateCode string
			RegionId     string
		}
	}
}

var Config = &config{}
var filePath = "config.toml"

// Init 初始化配置
func (c *config) Init() error {
	return configor.Load(Config, filePath)
}

// ENV 获取当前配置场景
func (c *config) ENV() string {
	return configor.ENV()
}

// Save 保存配置
func (c *config) Save() (err error) {
	var (
		file   *os.File
		buffer bytes.Buffer
	)
	if file, err = os.Create(filePath); err != nil {
		return
	}
	defer file.Close()
	err = toml.NewEncoder(&buffer).Encode(Config)
	if err != nil {
		return
	}
	_, err = file.Write(buffer.Bytes())
	return
}
