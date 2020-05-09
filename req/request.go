/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package req

import "strings"

const (
	PlatformID    string = "PlatformID"
	MiniProgramID        = "MiniProgramID"
	PayMchID             = "PayMchID"
)

// 分页基础入参
type PageParam struct {
	// id分页时使用，如无特殊说明不用此字段
	LastId   int `json:"last_id"`
	Page     int `json:"page" binding:"min=1"`
	PageSize int `json:"page_size" binding:"min=1,max=50"`
}

func (p *PageParam) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// 排序入参
type OrderBase struct {
	OrderColumn string // 排序列
	OrderClase  string // 排序方式，desc倒序， asc正序
}

func (p *OrderBase) GetSafeClase() string {
	p.OrderClase = strings.ToLower(p.OrderClase)
	switch p.OrderClase {
	case "asc":
		return "asc"
	case "desc":
		return "desc"
	default:
		return "desc"
	}
}
