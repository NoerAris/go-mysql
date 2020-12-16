package main

import (
	//For swagger
	"github.com/iris-contrib/swagger"
	"github.com/iris-contrib/swagger/swaggerFiles"
	_ "go-mysql/docs"

	"github.com/kataras/iris"
	"go-mysql/config"
	_ "go-mysql/controller"
	_ "go-mysql/controller/userController"
	"go-mysql/router"
	"os"
)

// @title CRUD MySql API
// @version 1.0
// @description This is a service for CRUD
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:9093
// @BasePath /

func main() {
	os.Unsetenv("log-files")
	app := iris.New()
	router.ProcessController(app)
	port :=config.GetValueFromConfig("", "app.port")

	config := &swagger.Config{
		URL: "http://localhost:" + port + "/swagger/doc.json", //The url pointing to API definition
	}
	// use swagger middleware to
	app.Get("/swagger/{any:path}", swagger.CustomWrapHandler(config, swaggerFiles.Handler))
	//Disable some environment
	//app.Get("/swagger/{any:path}", swagger.DisablingWrapHandler(swaggerFiles.Handler, "NAME_OF_ENV_VARIABLE"))
	app.Run(iris.Addr(":" + port))
}
