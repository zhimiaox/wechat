/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package utils

import (
	"fmt"
	"testing"
)

type B struct {
	A1 string
	B1 int
	G1 bool
}

type A struct {
	A1 string
}

func TestSuperConvert(t *testing.T) {
	a := A{
		A1: "123",
	}
	b := B{}
	SuperConvert(&a, &b)
	fmt.Printf("%#v", b)
}
