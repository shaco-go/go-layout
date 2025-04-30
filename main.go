package main

import (
	"flag"
	"forum/bootstrap"
	"forum/g"
	"forum/internal/server"
)

func main() {
	// 1.服务配置文件
	var configFile string
	flag.StringVar(&configFile, "c", "./config/development.yaml", "配置文件路径")
	flag.Parse()
	g.Conf = bootstrap.InitConfig(configFile)
	// 2.初始化日志
	bootstrap.InitLogger()
	g.DB = bootstrap.InitDB()
	g.Pool = bootstrap.InitPool()
	g.Redis = bootstrap.InitRedis()
	server.Run()
}
