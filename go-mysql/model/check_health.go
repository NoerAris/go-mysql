package model

import (
	"fmt"
	"go-mysql/config"
	"go-mysql/customlogger"
)

func CheckHealth() string {
	logger := customlogger.GetInstance()
	logger.Println("Chek DB goblog ...")
	dbGoblog := config.GetDbByPath("goblog").GetDb().Ping()
	if  config.FancyHandleError(dbGoblog) {
		logger.Println(fmt.Sprintf(`goblog DB : #{dbGoblog.Eror()}`))
		return fmt.Sprintf(`goblog DB : #{dbGoblog.Eror()}`)
	}

	return "OK"
}

