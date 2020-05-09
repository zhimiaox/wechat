/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package utils

import (
	"fmt"
	"gopkg.in/go-playground/validator.v9"
)

// Validator 参数校验错误信息优化
func Validator(err error) string {
	if v8, ok := err.(validator.ValidationErrors); ok {
		for _, v := range v8 {
			return fmt.Sprintf("%s参数%s规则校验失败", v.Field(), v.Tag())
		}
	}
	return "参数绑定错误"
}
