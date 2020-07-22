/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package main

import (
	"gitee.com/zhimiao/wechat/api"
	"gitee.com/zhimiao/wechat/common"
	"gitee.com/zhimiao/wechat/models"
)

// @title 纸喵 wechat API
// @version 1.0
// @description 纸喵软件系列之服务端
// @termsOfService http://zhimiao.org

// @contact.name API Support
// @contact.url http://tools.zhimiao.org
// @contact.email mail@xiaoliu.org

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host
// @basepath /
func main() {
	err := common.Config.Init()
	if err != nil {
		panic("配置加载失败")
	}
	models.Start()
	go api.Start()
	select {}
}
