/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package api

import (
	"fmt"
	"gitee.com/zhimiao/wechat/common"
	"gitee.com/zhimiao/wechat/common/utils"
	"gitee.com/zhimiao/wechat/models"
	"gitee.com/zhimiao/wechat/req"
	"gitee.com/zhimiao/wechat/resp"
	ws "gitee.com/zhimiao/wechat/wechat_service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type oplatformApi struct{}

var OPlatform = &oplatformApi{}

// @Summary 设施平台信息
// @Tags 第三方平台管理
// @Produce  json
// @Accept  json
// @Param body body req.PlatformParam true "入参集合"
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "操作成功","data": null}"
// @Router /oplatform/manage/Set [post]
func (m *oplatformApi) Set(c *gin.Context) {
	param := &req.PlatformParam{}
	if err := c.ShouldBind(param); err != nil {
		resp.NewApiResult(-4, utils.Validator(err)).Json(c)
		return
	}
	s := &models.Platform{}
	utils.SuperConvert(param, s)
	err := s.Set()
	if err != nil {
		resp.NewApiResult(-5, err.Error()).Json(c)
		return
	}
	resp.NewApiResult(1).Json(c)
}

// @Summary 获取平台列表
// @Tags 第三方平台管理
// @Produce  json
// @Accept  json
// @Param param query req.PageParam true "入参集合"
// @Success 200 {array} models.Platform ""
// @Router /oplatform/manage/Lists [get]
func (m *oplatformApi) Lists(c *gin.Context) {
	param := &req.PageParam{}
	if err := c.ShouldBind(param); err != nil {
		resp.NewApiResult(-4, utils.Validator(err)).Json(c)
		return
	}
	s := &models.Platform{}
	lists, rows := s.List(param.Offset(), param.PageSize)
	vo := make([]resp.PlatformListsVO, len(lists))
	for k, v := range lists {
		utils.SuperConvert(&v, &vo[k])
		vo[k].VPlatformState = ws.GetOPState(v.PlatformID)
	}
	resp.NewApiResult(1, "读取成功", resp.PageInfo{
		Page:      param.Page,
		PageSize:  param.PageSize,
		TotalSize: rows,
		Rows:      vo,
	}).Json(c)
}

// @Summary 刷新服务器、业务域名
// @Tags 小程序控制
// @Produce  json
// @Param MiniProgramID query string true "小程序APPID"
// @Param PlatformID path string true "第三方平台APPID"
// @Success 200 {object} resp.ApiResult ""
// @Router /oplatform/{PlatformID}/SetDomain [post]
func (m *oplatformApi) SetDomain(c *gin.Context) {
	miniProgramId := c.Query(req.MiniProgramID)
	if miniProgramId == "" {
		resp.NewApiResult(-4, "入参错误").Json(c)
		return
	}
	err := ws.SetMiniAppDomain(c.GetString(req.PlatformID), miniProgramId)
	if err != nil {
		resp.NewApiResult(-5, "设置失败").Json(c)
		logrus.Warn(err.Error())
		return
	}
	resp.NewApiResult(1, "操作成功").Json(c)
}

// @Summary 审核加急
// @Tags 小程序控制
// @Produce  json
// @Param body body req.SpeedUpAuditParam true "入参集合"
// @Param PlatformID path string true "第三方平台APPID"
// @Success 200 {object} resp.ApiResult ""
// @Router /oplatform/{PlatformID}/SpeedUpAudit [post]
func (m *oplatformApi) SpeedUpAudit(c *gin.Context) {
	param := &req.SpeedUpAuditParam{}
	if err := c.ShouldBind(param); err != nil {
		resp.NewApiResult(-4, utils.Validator(err)).Json(c)
		return
	}
	maService := ws.GetOpenMiniPrograms(c.GetString(req.PlatformID), param.MiniProgramID)
	if maService == nil {
		resp.NewApiResult(-5, "请求失败").Json(c)
		return
	}
	if err := maService.SpeedUpAudit(param.AuditID); err != nil {
		logrus.Warn("%s号小程序审核加急失败：%s", param.MiniProgramID, err.Error())
		resp.NewApiResult(-6, "加急失败").Json(c)
		return
	}
	resp.NewApiResult(1, "操作成功").Json(c)
}

// @Summary 发布已审核通过的小程序
// @Tags 小程序控制
// @Produce  json
// @Param MiniProgramID query string true "小程序APPID"
// @Param PlatformID path string true "第三方平台APPID"
// @Success 200 {object} resp.ApiResult ""
// @Router /oplatform/{PlatformID}/Release [post]
func (m *oplatformApi) Release(c *gin.Context) {
	miniProgramId := c.Query(req.MiniProgramID)
	if miniProgramId == "" {
		resp.NewApiResult(-4, "入参错误").Json(c)
		return
	}
	err := ws.Release(c.GetString(req.PlatformID), miniProgramId)
	if err != nil {
		resp.NewApiResult(-6, "发布失败").Json(c)
		return
	}
	resp.NewApiResult(1, "操作成功").Json(c)
}

// @Summary 提交代码
// @Tags 小程序控制
// @Produce  json
// @Accept  json
// @Param body body req.CommitCodeParam true "入参集合"
// @Param PlatformID path string true "第三方平台APPID"
// @Success 200 {object} resp.ApiResult ""
// @Router /oplatform/{PlatformID}/CommitCode [post]
func (m *oplatformApi) CommitCode(c *gin.Context) {
	param := &req.CommitCodeParam{}
	if err := c.ShouldBind(param); err != nil {
		resp.NewApiResult(-4, utils.Validator(err)).Json(c)
		return
	}
	PlatformID := c.GetString(req.PlatformID)
	tpl, err := ws.GetTemplateDetail(PlatformID, param.TemplateID)
	if err != nil {
		resp.NewApiResult(-5, "模板获取失败").Json(c)
		return
	}
	err = ws.CommitCode(PlatformID, param.MiniProgramID, tpl)
	if err != nil {
		resp.NewApiResult(-5, "设置失败").Json(c)
		logrus.Warn(err.Error())
		return
	}
	resp.NewApiResult(1, "操作成功").Json(c)
}

// @Summary 提交审核
// @Tags 小程序控制
// @Param MiniProgramID query string true "小程序APPID"
// @Param PlatformID path string true "第三方平台APPID"
// @Success 200 {object} resp.ApiResult ""
// @Router /oplatform/{PlatformID}/SubmitAudit [post]
func (m *oplatformApi) SubmitAudit(c *gin.Context) {
	miniProgramId := c.Query(req.MiniProgramID)
	if miniProgramId == "" {
		resp.NewApiResult(-4, "入参错误").Json(c)
		return
	}
	auditId, err := ws.SubmitAudit(c.GetString(req.PlatformID), miniProgramId)
	if err != nil {
		logrus.Warn(err.Error())
		resp.NewApiResult(-5, "请求失败").Json(c)
		return
	}
	resp.NewApiResult(1, "操作成功", auditId).Json(c)
}

// @Summary 获取小程序最后一次审核状态
// @Tags 小程序信息获取
// @Param MiniProgramID query string true "小程序APPID"
// @Param PlatformID path string true "第三方平台APPID"
// @Router /oplatform/{PlatformID}/GetLatestAuditStatus [get]
func (m *oplatformApi) GetLatestAuditStatus(c *gin.Context) {
	miniProgramId := c.Query(req.MiniProgramID)
	if miniProgramId == "" {
		resp.NewApiResult(-4, "入参错误").Json(c)
		return
	}
	maService := ws.GetOpenMiniPrograms(c.GetString(req.PlatformID), miniProgramId)
	if maService == nil {
		resp.NewApiResult(-5, "请求失败").Json(c)
		return
	}
	ret, err := maService.GetLatestAuditStatus()
	if err != nil {
		resp.NewApiResult(-5, "请求失败").Json(c)
		return
	}
	resp.NewApiResult(1, "操作成功", ret).Json(c)
}

// @Summary 撤销审核
// @Description 撤销审核 每天1次每月10次
// @Tags 小程序控制
// @Param MiniProgramID query string true "小程序APPID"
// @Param PlatformID path string true "第三方平台APPID"
// @Router /oplatform/{PlatformID}/UndoCodeAudit [post]
func (m *oplatformApi) UndoCodeAudit(c *gin.Context) {
	miniProgramId := c.Query(req.MiniProgramID)
	if miniProgramId == "" {
		resp.NewApiResult(-4, "入参错误").Json(c)
		return
	}
	maService := ws.GetOpenMiniPrograms(c.GetString(req.PlatformID), miniProgramId)
	if maService == nil {
		resp.NewApiResult(-5, "请求失败").Json(c)
		return
	}
	err := maService.UndoCodeAudit()
	if err != nil {
		resp.NewApiResult(-5, "撤回失败").Json(c)
		logrus.Warn("%s号小程序审核撤回失败：%s", miniProgramId, err.Error())
		return
	}
	// 更新小程序状态为已撤审
	audit := &models.MiniprogramAudit{
		AppID: miniProgramId,
		State: -1,
	}
	audit.ChangeState()
	resp.NewApiResult(1, "操作成功").Json(c)
}

// @Summary 获取体验二维码
// @Tags 小程序信息获取
// @Param MiniProgramID query string true "小程序APPID"
// @Param PlatformID path string true "第三方平台APPID"
// @Router /oplatform/{PlatformID}/GetTestQrcode [get]
func (m *oplatformApi) GetTestQrcode(c *gin.Context) {
	miniProgramId := c.Query(req.MiniProgramID)
	if miniProgramId == "" {
		resp.NewApiResult(-4, "入参错误").Json(c)
		return
	}
	maService := ws.GetOpenMiniPrograms(c.GetString(req.PlatformID), miniProgramId)
	if maService == nil {
		resp.NewApiResult(-5, "请求失败").Json(c)
		return
	}
	ret, err := maService.GetTestQrcode("")
	if err != nil {
		logrus.Warn("小程序体验码获取失败,%s", err.Error())
		c.Status(500)
		return
	}
	c.Data(http.StatusOK, "image/jpeg", ret)
}

// @Summary 获取已设置类目
// @Tags 小程序信息获取
// @Param MiniProgramID query string true "小程序APPID"
// @Param PlatformID path string true "第三方平台APPID"
// @Router /oplatform/{PlatformID}/GetCodeCategory [get]
func (m *oplatformApi) GetCodeCategory(c *gin.Context) {
	miniProgramId := c.Query(req.MiniProgramID)
	if miniProgramId == "" {
		resp.NewApiResult(-4, "入参错误").Json(c)
		return
	}
	maService := ws.GetOpenMiniPrograms(c.GetString(req.PlatformID), miniProgramId)
	if maService == nil {
		resp.NewApiResult(-5, "请求失败").Json(c)
		return
	}
	ret, err := maService.GetCategory()
	if err != nil {
		logrus.Warn(err.Error())
		resp.NewApiResult(-5, "请求失败").Json(c)
		return
	}
	resp.NewApiResult(1, "操作成功", ret).Json(c)
}

// @Summary 获取小程序线上页面列表
// @Tags 小程序信息获取
// @Param MiniProgramID query string true "小程序APPID"
// @Param PlatformID path string true "第三方平台APPID"
// @Router /oplatform/{PlatformID}/GetCodePageList [get]
func (m *oplatformApi) GetCodePageList(c *gin.Context) {
	miniProgramId := c.Query(req.MiniProgramID)
	if miniProgramId == "" {
		resp.NewApiResult(-4, "入参错误").Json(c)
		return
	}
	maService := ws.GetOpenMiniPrograms(c.GetString(req.PlatformID), miniProgramId)
	if maService == nil {
		resp.NewApiResult(-5, "请求失败").Json(c)
		return
	}
	ret, err := maService.GetCodePage()
	if err != nil {
		logrus.Warn(err.Error())
		resp.NewApiResult(-5, "请求失败").Json(c)
		return
	}
	resp.NewApiResult(1, "操作成功", ret).Json(c)
}

// @Summary 获取小程序授权信息
// @Tags 小程序信息获取
// @Param MiniProgramID query string true "小程序APPID"
// @Param PlatformID path string true "第三方平台APPID"
// @Router /oplatform/{PlatformID}/GetAuthorizerInfo [get]
func (m *oplatformApi) GetAuthorizerInfo(c *gin.Context) {
	miniProgramId := c.Query(req.MiniProgramID)
	if miniProgramId == "" {
		resp.NewApiResult(-4, "入参错误").Json(c)
		return
	}
	if miniProgramId == "" {
		resp.NewApiResult(-4, "入参错误-appid").Json(c)
		return
	}
	openService := ws.GetOpen(c.GetString(req.PlatformID))
	if openService == nil {
		resp.NewApiResult(-5, "请求失败").Json(c)
		return
	}
	info, base, err := openService.GetAuthrInfo(miniProgramId)
	if err != nil {
		logrus.Warn(err.Error())
		resp.NewApiResult(-5, "请求失败").Json(c)
		return
	}
	resp.NewApiResult(1, "操作成功", map[string]interface{}{
		"component_appid":  info,
		"authorizer_appid": base,
	}).Json(c)
}

// @Summary 删除模板
// @Tags 第三方平台模板
// @Param template_id query string true "模板编号"
// @Param PlatformID path string true "第三方平台APPID"
// @Router /oplatform/{PlatformID}/tpl/del [delete]
func (m *oplatformApi) TplDel(c *gin.Context) {
	openService := ws.GetOpen(c.GetString(req.PlatformID))
	if openService == nil {
		resp.NewApiResult(-5, "请求失败").Json(c)
		return
	}
	templateId, err := strconv.Atoi(c.Query("template_id"))
	if templateId == 0 || err != nil {
		resp.NewApiResult(-5, "模板ID无法识别").Json(c)
		return
	}
	if err := openService.DeleteTpl(templateId); err != nil {
		logrus.Warn(err.Error())
		resp.NewApiResult(-5, "请求失败").Json(c)
		return
	}
	resp.NewApiResult(1, "操作成功").Json(c)
}

// @Summary 获取模板列表
// @Tags 第三方平台模板
// @Param PlatformID path string true "第三方平台APPID"
// @Router /oplatform/{PlatformID}/tpl/list [get]
func (m *oplatformApi) TplList(c *gin.Context) {
	openService := ws.GetOpen(c.GetString(req.PlatformID))
	if openService == nil {
		resp.NewApiResult(-5, "请求失败").Json(c)
		return
	}
	ret, err := openService.TplList()
	if err != nil {
		logrus.Warn(err.Error())
		resp.NewApiResult(-5, "请求失败").Json(c)
		return
	}
	resp.NewApiResult(1, "操作成功", ret.TemplateList).Json(c)
}

// @Summary 将模板推送到所有自动升级小程序
// @Tags 第三方平台模板
// @Param template_id query string true "模板编号"
// @Param PlatformID path string true "第三方平台APPID"
// @Router /oplatform/{PlatformID}/tpl/pushToAuto [post]
func (m *oplatformApi) pushToAuto(c *gin.Context) {
	templateId, err := strconv.Atoi(c.Query("template_id"))
	if templateId == 0 || err != nil {
		resp.NewApiResult(-5, "模板ID无法识别").Json(c)
		return
	}
	err = ws.PushToAutoUpdateMiniApp(c.GetString(req.PlatformID), templateId)
	if err != nil {
		resp.NewApiResult(-5, err.Error()).Json(c)
		return
	}
	resp.NewApiResult(1, "操作成功").Json(c)
}

// @Summary 添加草稿到模板
// @Tags 第三方平台模板
// @Param draft_id query string true "草稿编号"
// @Param PlatformID path string true "第三方平台APPID"
// @Router /oplatform/{PlatformID}/tpl/add [post]
func (m *oplatformApi) TplAdd(c *gin.Context) {
	openService := ws.GetOpen(c.GetString(req.PlatformID))
	if openService == nil {
		resp.NewApiResult(-5, "请求失败").Json(c)
		return
	}
	draftId, err := strconv.Atoi(c.Query("draft_id"))
	if draftId == 0 || err != nil {
		resp.NewApiResult(-5, "草稿ID无法识别").Json(c)
		return
	}
	if err := openService.DeleteTpl(draftId); err != nil {
		logrus.Warn(err.Error())
		resp.NewApiResult(-5, "请求失败").Json(c)
		return
	}
	resp.NewApiResult(1, "操作成功").Json(c)
}

// @Summary 获取草稿列表
// @Tags 第三方平台模板
// @Param PlatformID path string true "第三方平台APPID"
// @Success 200 {object} string ""
// @Router /oplatform/{PlatformID}/tpl/draft [get]
func (m *oplatformApi) TplDraft(c *gin.Context) {
	openService := ws.GetOpen(c.GetString(req.PlatformID))
	if openService == nil {
		return
	}
	ret, err := openService.TplDraftList()
	if err != nil {
		logrus.Warn(err.Error())
		resp.NewApiResult(-5, "请求失败").Json(c)
		return
	}
	resp.NewApiResult(1, "操作成功", ret.DraftList).Json(c)
}

// @Summary 授权
// @Tags 授权回调
// @Param PlatformID path string true "第三方平台APPID"
// @Router /oplatform/{PlatformID}/auth [get]
func (m *oplatformApi) Auth(c *gin.Context) {
	openService := ws.GetOpen(c.GetString(req.PlatformID))
	if openService == nil {
		return
	}
	apiHost := common.Config.Server.ApiHost
	redirect := fmt.Sprintf("%s/oplatform/%s/redirect", apiHost, openService.AppID)
	authURL, err := openService.AuthURL(redirect, 2)
	if err != nil {
		c.String(500, "系统异常:"+err.Error())
		return
	}
	c.Header("Content-Type", "text/html")
	c.String(200, fmt.Sprintf("<script>location.href=\"%s\"</script>", authURL))
}

// @Summary 同步回调
// @Tags 授权回调
// @Param PlatformID path string true "第三方平台APPID"
// @Router /oplatform/{PlatformID}/redirect [get]
func (m *oplatformApi) Redirect(c *gin.Context) {
	platform := &models.Platform{
		PlatformID: c.GetString(req.PlatformID),
	}
	rows, _ := platform.GetByAppId()
	if rows == 0 {
		c.String(200, "平台ID错误！")
		return
	}
	c.Redirect(http.StatusFound, platform.AuthRedirectURL)
}

// @Summary 异步回调
// @Tags 授权回调
// @Param PlatformID path string true "第三方平台APPID"
// @Router /oplatform/{PlatformID}/notify [post]
func (m *oplatformApi) Notify(c *gin.Context) {
	opID := c.GetString(req.PlatformID)
	config := ws.GetOpenConfig(opID)
	wechatService := ws.GetWechatService(config)
	if wechatService == nil {
		return
	}
	// 传入request和responseWriter
	server := wechatService.GetServer(c.Request, c.Writer)
	server.SetDebug(true)
	// 设置接收消息的处理方法
	server.SetMessageHandler(ws.MessageHandler)
	// 处理消息接收以及回复
	err := server.Serve()
	// 记录回调错误
	vo := ws.GetOPState(opID)
	vo.Notify = ""
	if err != nil {
		vo.Notify = err.Error()
	}
	ws.RefreshOPState(opID, vo)
	if err != nil {
		logrus.Warn(err.Error())
		return
	}
	// 发送回复的消息
	server.Send()
}
