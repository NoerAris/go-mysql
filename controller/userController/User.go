package userController

import (
	//For swagger
	_ "go-mysql/docs"

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

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the input paylod
// @Tags Save User
// @Accept  json
// @Produce  json
// @Param user body user.User true "Create user"
// @Success 200 {object} user.ResponseUser
// @Router /user/saveUser [post]
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

// GetUserById godoc
// @Summary Get a user by id
// @Description Get user by id was input
// @Tags Get User
// @Produce  json
// @Param id path int true "Get user by id"
// @Success 200 {object} user.ResponseUser
// @Router /user/findOne/{id} [get]
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

// GetAllUser godoc
// @Summary Get all user
// @Description Get all user
// @Tags Get All User
// @Produce  json
// @Success 200 {object} user.ResponseUser
// @Router /user/findAll [get]
func getAllusers (ctx context.Context) {
	logger := customlogger.GetInstance()
	logger.Println("Starting get all users ...")

	users := user.GetAllUser()
	ctx.JSON(users)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update a user with the id and input paylod
// @Tags Update User
// @Accept  json
// @Produce  json
// @Param user body user.User true "Update user"
// @Param id path int true "Update user"
// @Success 200 {object} user.ResponseUser
// @Router /user/updateUser/{id} [put]
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

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by id was input
// @Tags Delete User
// @Produce	json
// @Param id path int true "Delete user by id"
// @Success 200 {object} user.ResponseUser
// @Router /user/deleteUser/{id} [delete]
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

//This sample query param in API => http://localhost:9093/user/findOne?id=10&name=Aris
// TestUrlPram godoc
// @Summary Tes url with query param
// @Description Tes url with query param
// @Tags Test Param
// @Produce  json
// @Param firstName query string true "First Name"
// @Param lastName query string true "Last Name"
// @Success 200 {string} string
// @Router /user/findOne [get]
func testUrlParam (ctx context.Context) {
	first := ctx.URLParam("firstName")
	last := ctx.URLParam("lastName")

	result := fmt.Sprintf("Data from parameter URL, first name : %v, last name : %v", first, last)
	fmt.Println(result)

	ctx.JSON(result)
}