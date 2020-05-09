/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package models

import (
	"testing"
)

func TestMiniprogramAudit_ChangeState(t *testing.T) {
	// 更新小程序状态为已撤审
	audit := &MiniprogramAudit{
		AppID:                "wxf79506a*********",
		OriginalID:           "miniApp.OriginalID",
		AuditID:              454545454,
		State:                1,
		Reason:               "",
		ScreenShot:           "",
		TemplateID:           123213,
		TemplateAppID:        "tplInfo.SourceMiniprogramAppid",
		TemplateAppName:      "tplInfo.SourceMiniprogram",
		TemplateAppDeveloper: "tplInfo.Developer",
		TemplateDesc:         "tplInfo.UserDesc",
		TemplateVersion:      "tplInfo.UserVersion",
	}
	audit2 := &MiniprogramAudit{
		OriginalID: audit.OriginalID,
	}
	audit.ChangeState()

	// audit2.State = -1
	// audit2.ChangeState()

	// audit2.State = 2
	// audit2.ChangeState()

	audit2.State = 3
	audit2.Reason = "失败原因"
	audit2.ChangeState()

}
