module gitee.com/zhimiao/wechat

go 1.13

require (
	gitee.com/zhimiao/wechat-sdk v1.1.1
	github.com/BurntSushi/toml v0.3.1
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/axgle/mahonia v0.0.0-20180208002826-3358181d7394
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/chenjiandongx/ginprom v0.0.0-20191022035802-6f3da3c84986
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.5.0
	github.com/go-openapi/jsonreference v0.19.3 // indirect
	github.com/go-openapi/spec v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.6 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/gomodule/redigo v2.0.1-0.20180627144507-2cd21d9966bf+incompatible
	github.com/influxdata/influxdb-client-go v0.1.5
	github.com/jinzhu/configor v1.1.1
	github.com/jinzhu/gorm v1.9.11
	github.com/json-iterator/go v1.1.8 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mailru/easyjson v0.7.0 // indirect
	github.com/mattn/go-isatty v0.0.11 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/prometheus/client_golang v1.2.1
	github.com/prometheus/client_model v0.0.0-20191202183732-d1d2010b5bee // indirect
	github.com/prometheus/procfs v0.0.8 // indirect
	github.com/siddontang/go v0.0.0-20180604090527-bdc77568d726
	github.com/sirupsen/logrus v1.4.2
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.3
	golang.org/x/crypto v0.0.0-20191219195013-becbf705a915
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553 // indirect
	golang.org/x/sys v0.0.0-20191219235734-af0d71d358ab // indirect
	golang.org/x/tools v0.0.0-20191219230827-5e752206af05 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/yaml.v2 v2.2.7 // indirect
)
// 鉴于此工程与sdk工程同步开发，变动频繁，建议同时下载两个工程运行
replace gitee.com/zhimiao/wechat-sdk v1.1.1 => ../wechat-sdk
