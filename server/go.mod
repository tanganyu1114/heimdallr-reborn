module gin-vue-admin

go 1.14

require (
	github.com/ClessLi/bifrost v1.0.7
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/casbin/casbin v1.9.1
	github.com/casbin/casbin/v2 v2.11.0
	github.com/casbin/gorm-adapter/v3 v3.0.2
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/fsnotify/fsnotify v1.5.1
	github.com/fvbock/endless v0.0.0-20170109170031-447134032cb6
	github.com/gin-gonic/gin v1.7.4
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gookit/color v1.3.1
	github.com/gorilla/websocket v1.4.2
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/jordan-wright/email v0.0.0-20200824153738-3f5bafa1cd84
	github.com/lestrrat-go/file-rotatelogs v2.3.0+incompatible
	github.com/lestrrat-go/strftime v1.0.3 // indirect
	github.com/marmotedu/errors v1.0.2
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mojocn/base64Captcha v1.3.1
	github.com/pkg/errors v0.9.1
	github.com/qiniu/api.v7/v7 v7.4.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
	github.com/shirou/gopsutil v2.20.8+incompatible
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.1
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.7
	github.com/tebeka/strftime v0.1.3 // indirect
	github.com/unrolled/secure v1.0.7
	go.uber.org/zap v1.19.1
	google.golang.org/grpc v1.43.0
	gorm.io/driver/mysql v1.1.2
	gorm.io/gorm v1.22.4
)

replace (
	github.com/ClessLi/bifrost => github.com/ClessLi/bifrost v1.0.8-0.20220511061302-192174cd668b
	github.com/casbin/gorm-adapter/v3 => github.com/casbin/gorm-adapter/v3 v3.0.2
	github.com/fsnotify/fsnotify => github.com/fsnotify/fsnotify v1.4.9
	github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.6.3
	github.com/go-openapi/swag => github.com/go-openapi/swag v0.19.8
	github.com/go-redis/redis => github.com/go-redis/redis v6.15.7+incompatible
	github.com/go-sql-driver/mysql => github.com/go-sql-driver/mysql v1.5.0
	github.com/robfig/cron/v3 => github.com/robfig/cron/v3 v3.0.0
	github.com/satori/go.uuid => github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra => github.com/spf13/cobra v1.1.1
	github.com/spf13/viper => github.com/spf13/viper v1.7.0
	go.uber.org/zap => go.uber.org/zap v1.13.0
	gorm.io/driver/mysql => gorm.io/driver/mysql v0.3.0
	gorm.io/gorm => gorm.io/gorm v1.20.5
)
