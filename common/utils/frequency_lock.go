/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

// 频率锁
package utils

import (
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type LockItem struct {
	Key      string
	LifeSpan time.Duration // 生命周期
	CreateOn time.Time     // 创建时间
}

type LockTable struct {
	sync.RWMutex
	CleanerDuraction time.Duration       // 触发定时清理器的时间
	Cleaner          *time.Timer         // 定时清理器
	Items            map[string]LockItem // 子集
}

func NewLockTable() *LockTable {
	return &LockTable{
		Items: make(map[string]LockItem),
	}
}

func (this *LockTable) IsLock(key string, lock_time time.Duration) bool {
	this.Lock()
	item := LockItem{
		Key:      key,
		LifeSpan: lock_time,
		CreateOn: time.Now(),
	}
	if item, ok := this.Items[key]; ok {
		this.Unlock()
		if item.LifeSpan-time.Now().Sub(item.CreateOn) < 0 {
			this.cleanerCheck()
			return false
		}
		logrus.Info("[LockTable] %s is limited\n", key)
		return true
	}
	this.Items[key] = item
	cleannerDuraction := this.CleanerDuraction
	this.Unlock()
	logrus.Info("[LockTable] add %s to table\n", key)
	if cleannerDuraction == 0 {
		this.cleanerCheck()
	}
	return false
}

func (this *LockTable) cleanerCheck() {
	this.Lock()
	defer this.Unlock()
	logrus.Info("[LockTable] start timer cleaner Duraction after %.2f s\n", this.CleanerDuraction.Seconds())
	if this.Cleaner != nil {
		this.Cleaner.Stop()
	}

	// 遍历当前限制的key, 遇到过期的将其删掉
	// 其余的则从中找到最近一个将要过期的key并且将它还有多少时间过期作为下一次清理任务的定时时间
	now := time.Now()
	smallestDuracton := 0 * time.Second
	for key, item := range this.Items {
		lifeSpan := item.LifeSpan
		createOn := item.CreateOn
		if now.Sub(createOn) >= lifeSpan {
			logrus.Info("[LockTable] delete key %s", key)
			delete(this.Items, key)
		} else {
			if smallestDuracton == 0 || lifeSpan-now.Sub(createOn) < smallestDuracton {
				smallestDuracton = lifeSpan - now.Sub(createOn)
			}
		}
	}

	this.CleanerDuraction = smallestDuracton
	// 将最近一个将要过期的key距离现在的时间作为启动清理任务的定时时间
	if this.CleanerDuraction > 0 {
		fn := func() {
			go this.cleanerCheck()
		}
		this.Cleaner = time.AfterFunc(this.CleanerDuraction, fn)
	}
}
