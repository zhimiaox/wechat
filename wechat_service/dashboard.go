/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package wechat_service

import (
	"gitee.com/zhimiao/wechat/resp"
	"sync"
)

var (
	opState *OpStateMap // 平台状态存储
)

type OpStateMap struct {
	sync.RWMutex
	Map map[string]resp.PlatformStateVO
}

func init() {
	opState = new(OpStateMap)
	opState.Map = make(map[string]resp.PlatformStateVO)
}

// GetOPState 获取平台状态
func GetOPState(k string) (vo resp.PlatformStateVO) {
	opState.RLock()
	vo = opState.Map[k]
	opState.RUnlock()
	return
}

// RefreshOPState 刷新平台状态
func RefreshOPState(k string, v resp.PlatformStateVO) {
	opState.Lock()
	opState.Map[k] = v
	opState.Unlock()
}
