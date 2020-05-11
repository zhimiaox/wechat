# 纸喵中间件 - 微信

[![giteego build-test](https://gitee.com/zhimiao/wechat/badge/giteego.svg?name=build-test&id=8258)](https://gitee.com/zhimiao/dashboard/projects/zhimiao/wechat/giteegos/8258?branch=master)

#### 介绍

集成微信第三方平台、小程序、支付等功能api,没有添加权限校验,通过简易修改即可整合到sass平台系统之中
此工程配套有对应的ui工程，vue组件
当前版本主打功能为第三方平台管控小程序系统，主要完善了：

    1.多个第三方平台添加管理
    2.第三方平台代小程序代码模板管理
    3.小程序授权给第三方平台
    4.第三方平台代小程序提交代码、提交审核、撤销审核、加急审核、发布
    5.小程序基础信息配置
    6.sessionkey信息开放，方便第三方整合自己的原生功能(待开发..)

#### 开发说明

> 基础要求

此工程提供了vue的前端组件，仓库传送门 [gitee.com/zhimiao/wechat-vue](https://gitee.com/zhimiao/wechat-vue)

鉴于此工程与[wechat-sdk](https://gitee.com/zhimiao/wechat-sdk)工程同步开发，变动频繁，建议同时下载两个工程运行

工程依赖MySQL数据库、Redis缓存请预先配置

> 数据表反序列化

```shell script
go get -u github.com/xxjwxc/gormt
gormt -H "localhost" -u "root" -p "root" -d "wechat_platform" -o "./tools/models"
```

> 文档自动生成

```shell script
go get -u github.com/swaggo/swag/cmd/swag
swag init
```

> Redis字典

| key | 备注 |
|:------|:-------|
| qy_access_token_${小程序APPID} | 小程序token |
| authorizer_access_token_${小程序APPID} | 代小程序accesstoken |
| component_access_token_${平台APPID} | 代小程序accesstoken |
| component_verify_ticket_${平台APPID} | 第三方平台票据 |

> 数据字典

[wechat_platform.sql](./wechat_platform.sql)

#### SDK挂件

[![纸喵软件/wechat-sdk](https://gitee.com/zhimiao/wechat-sdk/widgets/widget_card.svg?colors=4183c4,ffffff,ffffff,e3e9ed,666666,9b9b9b)](https://gitee.com/zhimiao/wechat-sdk)