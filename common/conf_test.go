/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package common

import (
	"testing"
)

func TestSave(t *testing.T) {
	Config.Init()
	Config.Save()
}
