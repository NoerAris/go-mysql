package controller

import (
	"github.com/kataras/iris/context"
	"go-mysql/customlogger"
	"go-mysql/model"
	"go-mysql/router"
)

func init () {
	c := router.CreateNewControllerInstance("clear", "/health")
	c.Get("", checkHealth)
}

func checkHealth (ctx context.Context) {
	logger := customlogger.GetInstance()
	logger.Println("Check health API's running ...")
	hasil := model.CheckHealth()
	ctx.Text(hasil)
}