package userController

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/context"
	"go-mysql/config"
	"go-mysql/customlogger"
	"go-mysql/model/user"
	"go-mysql/router"
	"strconv"
)

func init () {
	c := router.CreateNewControllerInstance("clear", "/user")
	c.Post("/saveUser", saveUser)
	c.Get("/findOne/{id}", getUserById)
	c.Get("/findAll", getAllusers)
	c.Put("/updateUser/{id}", updateUser)
	c.Delete("/deleteUser/{id}", deleteUser)
	c.Get("/findOne", testUrlParam)
}

func saveUser(ctx context.Context) {
	var userModel user.User
	logger := customlogger.GetInstance()
	err := json.NewDecoder(ctx.Request().Body).Decode(&userModel)

	logger.Println("Starting insert user ...")

	if err != nil {
		logger.Printf("Unable to decode the request body.  %v", err)
		response := user.ResponseUser{
			Message: fmt.Sprintf("Unable to decode the request body.  %v", err),
			Data: nil,
		}
		ctx.JSON(response)
		return
	}

	usr := user.InsertUser(userModel)
	ctx.JSON(usr)
}

func getUserById (ctx context.Context) {
	param := ctx.Params().Get("id")
	id, err := strconv.Atoi(param)

	logger := customlogger.GetInstance()
	logger.Println("Get data user by id : " + param)

	if err != nil {
		logger.Println(fmt.Sprintf("Unable to convert the string %v into int. %v", id, err))
		response := user.ResponseUser{
			Message: fmt.Sprintf("Unable to convert the string %v into int. %v", id, err),
			Data: nil,
		}
		ctx.JSON(response)
		return
	}

	dataUser := user.GetUser(int64(id))
	ctx.JSON(dataUser)
}

func getAllusers (ctx context.Context) {
	logger := customlogger.GetInstance()
	logger.Println("Starting get all users ...")

	users := user.GetAllUser()
	ctx.JSON(users)
}

func updateUser (ctx context.Context) {
	var dataUser user.User
	param := ctx.Params().Get("id")
	id, err := strconv.Atoi(param)
	errData := json.NewDecoder(ctx.Request().Body).Decode(&dataUser)

	logger := customlogger.GetInstance()
	logger.Println("Starting update user with id : " + param)

	if config.FancyHandleError(errData) {
		logger.Println("Unable to decode request body, " + errData.Error())
		response := user.ResponseUser{
			Message: "Unable to decode request body, " + errData.Error(),
			Data: nil,
		}
		ctx.JSON(response)
		return
	}

	if config.FancyHandleError(err) {
		logger.Printf("Unable to convert the string %v into int. %v", id, err)
		response := user.ResponseUser{
			Message: fmt.Sprintf("Unable to convert the string %v into int. %v", id, err),
			Data: nil,
		}
		ctx.JSON(response)
		return
	}

	result := user.UpdateUser(int64(id), dataUser)
	ctx.JSON(result)
}

func deleteUser (ctx context.Context) {
	logger := customlogger.GetInstance()
	param := ctx.Params().Get("id")
	logger.Println("Starting delete user with id : " + param)

	id, err := strconv.Atoi(param)
	if config.FancyHandleError(err) {
		logger.Printf("Unable to convert string %v to int. %v", id, err)
		response := user.ResponseUser{
			Message: fmt.Sprintf("Unable to convert string %v to int. %v", id, err),
			Data: nil,
		}
		ctx.JSON(response)
		return
	}

	result := user.DeleteUser(int64(id))
	ctx.JSON(result)
}

//This sample API => http://localhost:9093/user/findOne?id=10&name=Aris
func testUrlParam (ctx context.Context) {
	id := ctx.URLParam("id")
	name := ctx.URLParam("name")

	result := fmt.Sprintf("Data from parameter URL, id : %v, name : %v", id, name)
	fmt.Println(result)

	ctx.JSON(result)
}