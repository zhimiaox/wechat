/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package utils

import (
	"fmt"
	"runtime"
	"sort"
)

func IdsFitter(ids []int) []int {
	sort.Ints(ids)
	var newIds []int
	var lastId int
	for k, v := range ids {
		if k == 0 {
			lastId = v
			newIds = append(newIds, v)
		}
		if k > 0 && v != lastId {
			lastId = v
			newIds = append(newIds, v)
		}
	}
	return newIds
}

// PanicToError Panic转换为error
func PanicToError(f func()) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf(PanicTrace(e))
		}
	}()
	f()
	return
}

// PanicTrace panic调用链跟踪
func PanicTrace(err interface{}) string {
	stackBuf := make([]byte, 4096)
	n := runtime.Stack(stackBuf, false)

	return fmt.Sprintf("panic: %v %s", err, stackBuf[:n])
}
