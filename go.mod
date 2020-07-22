module gitee.com/zhimiao/wechat

go 1.14

require (
	gitee.com/zhimiao/wechat-sdk v1.1.2
	github.com/BurntSushi/toml v0.3.1
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/axgle/mahonia v0.0.0-20180208002826-3358181d7394
	github.com/chenjiandongx/ginprom v0.0.0-20200410120253-7cfb22707fa6
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/gomodule/redigo v2.0.1-0.20180627144507-2cd21d9966bf+incompatible
	github.com/influxdata/influxdb-client-go v1.1.0
	github.com/jinzhu/configor v1.2.0
	github.com/jinzhu/gorm v1.9.12
	github.com/prometheus/client_golang v1.6.0
	github.com/siddontang/go v0.0.0-20180604090527-bdc77568d726
	github.com/sirupsen/logrus v1.6.0
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.5
	golang.org/x/crypto v0.0.0-20200429183012-4b2356b1ed79
)

// 鉴于此工程与sdk工程同步开发，变动频繁，建议同时下载两个工程运行
replace gitee.com/zhimiao/wechat-sdk v1.1.2 => ../wechat-sdk
