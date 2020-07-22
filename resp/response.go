/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package resp

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type ApiResult struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type PageInfo struct {
	Page      int
	PageSize  int
	TotalSize int
	Rows      interface{}
}

func NewApiResult(arg ...interface{}) *ApiResult {
	result := &ApiResult{
		Code: 1,
		Msg:  "操作成功",
	}
	for k, v := range arg {
		if k == 0 {
			if v1, ok := v.(int); ok {
				result.setCode(v1)
			}
		}
		if k == 1 {
			if v2, ok := v.(string); ok {
				result.setMsg(v2)
			}
		}
		if k == 2 && v != nil {
			result.setData(v)
		}
	}
	return result
}

func (r *ApiResult) setData(data interface{}) *ApiResult {
	r.Data = data
	return r
}

func (r *ApiResult) setMsg(msg string) *ApiResult {
	r.Msg = msg
	return r
}

func (r *ApiResult) setCode(code int) *ApiResult {
	r.Code = code
	return r
}

func (r *ApiResult) Json(c *gin.Context) {
	c.JSON(200, r)
}

type JsonTimeUnix time.Time
type JsonTimeDate time.Time

func (t JsonTimeDate) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf(`"%s"`, time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}
func (t JsonTimeUnix) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf(`"%d"`, time.Time(t).Unix())
	return []byte(stamp), nil
}
