/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package utils

// Validator 参数校验错误信息优化
func Validator(err error) string {
	return "参数绑定错误" + err.Error()
}
