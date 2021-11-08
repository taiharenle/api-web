package main

import (
	"api-web/conf"
	"api-web/server"
)

func main() {
	// 加载配置
	conf.Init()

	// 装载路由
	r := server.NewRouter()

	// 启动！
	r.Run(":3000")
}
