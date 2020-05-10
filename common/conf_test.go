/*
 * 纸喵软件
 * Copyright (c) 2017~2020 zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/5/9 下午11:10
 * LastModified: 2020/5/9 下午11:07
 */

package common

import (
	"testing"
)

func TestSave(t *testing.T) {
	Config.Init()
	Config.Save()
}
