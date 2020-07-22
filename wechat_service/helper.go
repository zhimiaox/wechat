/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package wechat_service

import (
	"encoding/json"
	"fmt"
	"gitee.com/zhimiao/wechat-sdk/open"
	"gitee.com/zhimiao/wechat/common"
	"gitee.com/zhimiao/wechat/models"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

// SetMiniAppDomain 修改小程序域名配置
func SetMiniAppDomain(PlatformID, miniProgramId string) (err error) {
	maService := GetOpenMiniPrograms(PlatformID, miniProgramId)
	if maService == nil {
		err = fmt.Errorf("请求失败")
		return
	}
	s := &models.Platform{
		PlatformID: PlatformID,
	}
	if rows, e2 := s.GetByAppId(); rows == 0 || e2 != nil {
		err = fmt.Errorf("平台信息获取失败")
		return
	}
	var e1, e2 error
	serverDomain := strings.Split(s.ServerDomain, ",")
	if len(serverDomain) > 0 {
		serverDomainParam := open.ModifyDomainParam{
			Action:          open.ActionSet,
			RequestDomain:   serverDomain,
			WSRequestDomain: serverDomain,
			UploadDomain:    serverDomain,
			DownloadDomain:  serverDomain,
		}
		e1 = maService.ModifyDomain(serverDomainParam)
		if e1 != nil {
			err = fmt.Errorf("服务器域名设置失败, %s", e1.Error())
		}
	}
	bizDomain := strings.Split(s.BizDomain, ",")
	if len(bizDomain) > 0 {
		bizDomainParam := open.SetWebViewDomainURLParam{
			Action:        open.ActionSet,
			WebViewDomain: bizDomain,
		}
		e2 = maService.SetWebViewDomain(bizDomainParam)
		if e2 != nil {
			err = fmt.Errorf("业务域名设置失败, %s", e2.Error())
		}
	}
	return
}

// GetTemplateDetail 根据模板ID获取模板详情
func GetTemplateDetail(PlatformID string, TemplateID int) (tpl open.TplDetail, err error) {
	cacheId := fmt.Sprintf("TemplateCache:%s", PlatformID)
	redis := models.Redis
	data := redis.HGet(cacheId, strconv.Itoa(TemplateID))
	var ok bool
	if tpl, ok = data.(open.TplDetail); ok {
		return
	}
	openService := GetOpen(PlatformID)
	if openService == nil {
		err = fmt.Errorf("开放平台初始化失败")
		return
	}
	ret, err := openService.TplList()
	if err != nil {
		return
	}
	for _, v := range ret.TemplateList {
		if v.TemplateID == TemplateID {
			redis.HSet(cacheId, strconv.Itoa(TemplateID), v)
			tpl = v
			return
		}
	}
	err = fmt.Errorf("未找到对应模板信息")
	return
}

// CommitCode 提交代码
func CommitCode(PlatformID string, MiniProgramID string, template open.TplDetail) (err error) {
	maService := GetOpenMiniPrograms(PlatformID, MiniProgramID)
	if maService == nil {
		err = fmt.Errorf("小程序句柄获取失败")
		return
	}
	miniApp := &models.Miniprogram{
		AppID: MiniProgramID,
	}
	rows, err := miniApp.GetByAppID()
	if err != nil || rows == 0 {
		err = fmt.Errorf("小程序配置信息获取失败")
		return
	}
	var extConfig = make(map[string]string)
	if miniApp.ExtConfig != "" {
		err = json.Unmarshal([]byte(miniApp.ExtConfig), &extConfig)
		if err != nil {
			err = fmt.Errorf("小程序配置信息解析失败")
			return
		}
	}
	commitParam := open.CommitParam{
		TemplateID: template.TemplateID,
		Ext: open.CommitParamExt{
			Ext: extConfig,
		},
		UserVersion: template.UserVersion,
		UserDesc:    template.UserDesc,
	}
	err = maService.Commit(commitParam)
	if err != nil {
		logrus.Warn(err.Error())
		err = fmt.Errorf("提交代码失败")
		return
	}
	// 更新小程序配置，刷新当前模板编号
	miniApp = &models.Miniprogram{
		AppID:          MiniProgramID,
		NowTemplateID:  template.TemplateID,
		TemplateListen: template.SourceMiniprogramAppid,
	}
	miniApp.UpdateByAppID()
	err = SetMiniAppDomain(PlatformID, MiniProgramID)
	if err != nil {
		logrus.Warn("尝试注入域名失败，%s", err.Error())
	}
	return
}

// SubmitAudit 提审
func SubmitAudit(PlatformID, miniProgramId string) (auditId uint64, err error) {
	maService := GetOpenMiniPrograms(PlatformID, miniProgramId)
	if maService == nil {
		err = fmt.Errorf("小程序句柄获取失败")
		return
	}
	// 获取当前小程序的审核状态
	lastAudit, err := maService.GetLatestAuditStatus()
	if err != nil {
		logrus.Warn("%s号小程序%s失败：%s", miniProgramId, "查询当前审核状态", err.Error())
		err = fmt.Errorf("查询当前审核状态失败，无法判断是否能够提审")
		return
	}
	if lastAudit.Status == 2 {
		err = fmt.Errorf("当前有正在审核的版本，请勿多次提交")
		return
	}
	// 查询是否还有提审额度
	quota, err := maService.QueryQuota()
	if err != nil {
		logrus.Warn("quota查询失败，%s", err.Error())
		err = fmt.Errorf("quota查询失败，无法判断是否能够提审")
		return
	}
	if quota.Rest == 0 {
		err = fmt.Errorf("当月分配提审次数%d已用尽", quota.Limit)
		return
	}
	// 获取类目
	auditCateList, err := maService.GetAuditCategory()
	if err != nil || len(auditCateList) == 0 {
		logrus.Warn("%s号小程序%s失败：%s", miniProgramId, "类目获取", err.Error())
		err = fmt.Errorf("类目获取失败，请至小程序后台添加经营类目")
		return
	}
	auditCate := auditCateList[0]
	// 获取代码页面
	pages, err := maService.GetCodePage()
	if err != nil || len(pages.PageList) == 0 {
		logrus.Warn("%s号小程序%s失败：%s", miniProgramId, "页面获取", err.Error())
		err = fmt.Errorf("待审小程序页面获取失败")
		return
	}
	// 组装审核项
	item := open.SubmitAuditItem{
		Address: pages.PageList[0],
		// Tag:         "",
		FirstClass:  auditCate.FirstClass,
		SecondClass: auditCate.SecondClass,
		ThirdClass:  auditCate.ThirdClass,
		FirstID:     auditCate.FirstID,
		SecondID:    auditCate.SecondID,
		ThirdID:     auditCate.ThirdID,
		Title:       "首页",
	}
	// 合成最后审核信息
	submitParam := open.SubmitAuditParam{
		ItemList:      []open.SubmitAuditItem{},
		FeedbackInfo:  "",
		FeedbackStuff: "",
	}
	submitParam.ItemList = append(submitParam.ItemList, item)
	// 获取小程序信息
	miniApp := &models.Miniprogram{
		AppID: miniProgramId,
	}
	rows, err := miniApp.GetByAppID()
	if err != nil || rows == 0 {
		err = fmt.Errorf("小程序信息获取失败")
		return
	}
	// 获取模板信息
	tplInfo, err := GetTemplateDetail(PlatformID, miniApp.NowTemplateID)
	if err != nil {
		return
	}
	// 发起审核
	auditId, err = maService.SubmitAudit(submitParam)
	if err != nil {
		logrus.Warn("%s号小程序提审失败：%s", miniProgramId, err.Error())
		err = fmt.Errorf("提审失败")
		return
	}
	// 记录审核单据
	auditModel := &models.MiniprogramAudit{
		AppID:                miniProgramId,
		OriginalID:           miniApp.OriginalID,
		AuditID:              auditId,
		State:                1,
		Reason:               "",
		ScreenShot:           "",
		TemplateID:           miniApp.NowTemplateID,
		TemplateAppID:        tplInfo.SourceMiniprogramAppid,
		TemplateAppName:      tplInfo.SourceMiniprogram,
		TemplateAppDeveloper: tplInfo.Developer,
		TemplateDesc:         tplInfo.UserDesc,
		TemplateVersion:      tplInfo.UserVersion,
	}
	auditModel.ChangeState()
	return
}

// Release 发布
func Release(PlatformID, miniProgramId string) (err error) {
	maService := GetOpenMiniPrograms(PlatformID, miniProgramId)
	if maService == nil {
		err = fmt.Errorf("小程序句柄获取失败")
		return
	}
	lastAudit, err := maService.GetLatestAuditStatus()
	if err != nil || lastAudit.Status != 0 {
		err = fmt.Errorf("没有小程序可发布")
		return
	}
	if err = maService.Release(); err != nil {
		logrus.Warn("%s号小程序发布失败：%s", miniProgramId, err.Error())
		err = fmt.Errorf("发布失败")
		return
	}
	// 查询审核数据
	auditModel := models.MiniprogramAudit{
		AppID:   miniProgramId,
		AuditID: lastAudit.AuditID,
	}
	auditModel.GetBySelectKey()
	// 更正小程序发布状态
	miniApp := &models.Miniprogram{
		AppID: miniProgramId,
		// Version: auditModel.TemplateVersion,
		State: 5,
	}
	miniApp.UpdateByAppID()
	return
}

// PushToAutoUpdateMiniApp 推送至自动升级小程序
func PushToAutoUpdateMiniApp(PlatformID string, templateId int) (err error) {
	template, err := GetTemplateDetail(PlatformID, templateId)
	if err != nil {
		return
	}
	miniApp := &models.Miniprogram{
		PlatformID:     PlatformID,
		NowTemplateID:  templateId,
		TemplateListen: template.SourceMiniprogramAppid,
	}
	list := miniApp.ListAutoAudit()
	if len(list) == 0 {
		err = fmt.Errorf("没有检测到自动升级的小程序")
		return
	}
	go func() {
		chanelLock := make(chan int, common.Config.App.MaxAuditProgress)
		for _, v := range list {
			chanelLock <- 1
			go func() {
				defer func() { <-chanelLock }()
				err := CommitCode(PlatformID, v.AppID, template)
				if err != nil {
					logrus.Warn("[%s]自动升级失败-提包，%s", v.AppID, err.Error())
					return
				}
				auditId, err := SubmitAudit(PlatformID, v.AppID)
				if err != nil {
					logrus.Warn("[%s]自动升级失败-提审，%s", v.AppID, err.Error())
					return
				}
				logrus.Info("[%s]自动升级成功，审核单号: %d", v.AppID, auditId)
			}()
		}
	}()
	return
}
