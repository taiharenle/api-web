package conf

import (
	"api-web/cache"
	"api-web/model"
	"api-web/util"
	"github.com/joho/godotenv"
	"os"
)

// Init 初始化配置项
func Init() {
	// 读取.env配置
	godotenv.Load()

	// 设置日志级别
	util.BuildLogger(os.Getenv("LOG_LEVEL"))

	// 连接数据库
	model.Database(os.Getenv("MYSQL_DSN"))

	// 连接Redis
	cache.Redis()
}
