package main

import (
	"github.com/kataras/iris"
	"go-mysql/config"
	_ "go-mysql/controller"
	_ "go-mysql/controller/userController"
	"go-mysql/customlogger"
	"go-mysql/router"
	"os"
)

func main() {
	os.Unsetenv("log-files")
	app := iris.New()
	router.ProcessController(app)
	port :=config.GetValueFromConfig("", "app.port")

	logger :=customlogger.GetInstance()
	logger.Println("Starting web server in port : " + port)
	app.Run(iris.Addr(":" + port))
}
