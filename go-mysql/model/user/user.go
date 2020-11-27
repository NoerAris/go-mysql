package user

import (
	"fmt"
	"go-mysql/config"
	"go-mysql/connection"
	"go-mysql/customlogger"
	"strconv"
)

//Order variables of struct must same with in table users and call body in postman
/*
== samples in post ==
{
    "name": "Ferdian",
    "age":29,
    "location":"Indonesia"
}
 */

type User struct {
	ID			int64 `json:"id,omitempty"`
	Name		string `json:"name"`
	Age			int8 `json:"age,omitempty"`
	Location 	string `json:"location"`
}

type ResponseUser struct {
	Message		string `json:"message"`
	Data		[]User    `json:"data"`
}

func InsertUser(usr User) ResponseUser {
	logger :=customlogger.GetInstance()
	logger.Println("Starting proccess insert user ...")

	//dbgoblog := config.GetDbByPath("goblog").GetDb()
	//dbgoblog.SetMaxIdleConns(0)
	//dbgoblog.SetConnMaxLifetime(300 * time.Second)

	var userModel []User
	dbgoblog := connection.GetGoblogConn()
	sqlStat, err := dbgoblog.Prepare("INSERT INTO users (name, age, location) VALUES (?, ?, ?)")
	defer sqlStat.Close()

	//Exec executes a prepared statement with the given arguments and returns a Result summarizing the effect of the statement.
	res, er := sqlStat.Exec(usr.Name, usr.Age, usr.Location)

	if err != nil || config.FancyHandleError(er) {
		logger.Printf("Unable to execute the query. %v, %v", er, err)
		response := ResponseUser{
			Message: fmt.Sprintf("Unable to execute the query. %v, %v", er, err),
			Data: nil,
		}
		sqlStat.Close()
		return response
	}

	//Get id from data was inserted
	var s int64
	s, _ = res.LastInsertId()
	str := strconv.FormatInt(s, 10)
	usr.ID = s

	response := ResponseUser{
		Message: fmt.Sprintf("Success inserted a single record with id : %v", str),
		Data: append(userModel, usr),
	}

	logger.Printf("Success inserted a single record with id : %v", str)
	return response
}

func GetUser(id int64) ResponseUser {
	logger := customlogger.GetInstance()
	logger.Println(fmt.Sprintf("Started get data user by id : %v", id))
	dbgoblog := connection.GetGoblogConn()

	var userModel User
	var userData []User
	sqlStat := `SELECT * FROM users where userid = ?`

	//QueryRow executes a query that is expected to return at most one row
	result := dbgoblog.QueryRow(sqlStat, id)

	//Scan copies the columns from the matched row into the values pointed at by dest
	err := result.Scan(&userModel.ID, &userModel.Name, &userModel.Age, &userModel.Location)
	if config.FancyHandleError(err) {
		logger.Println("Scan query select result : Error " + err.Error())
		response := ResponseUser{
			Message: "Scan query select result : Error " + err.Error(),
			Data: nil,
		}
		return response
	}

	response := ResponseUser{
		Message: "Success",
		Data: append(userData, userModel),
	}
	return response
}

func GetAllUser () ResponseUser {
	logger := customlogger.GetInstance()
	dbgoblog := connection.GetGoblogConn()

	var userModel []User
	sqlStat := `SELECT * FROM users`

	//Query executes a query that returns rows, typically a SELECT
	rows, err := dbgoblog.Query(sqlStat)
	defer rows.Close()

	if config.FancyHandleError(err) {
		logger.Println(fmt.Sprintf("Unable to execute the query. %v", err))
		rows.Close()
		response := ResponseUser{
			Message: err.Error(),
			Data: nil,
		}
		return response
	}

	for rows.Next() {
		var userData User
		er := rows.Scan(&userData.ID, &userData.Name, &userData.Age, &userData.Location)
		if er != nil {
			logger.Println(fmt.Sprintf(`Unable to scan the row %v`, er))
		} else {
			userModel = append(userModel, userData)
		}
	}
	response := ResponseUser{
		Message: "Success",
		Data: userModel,
	}
	return response
}

func UpdateUser (id int64, usr User) ResponseUser {
	logger := customlogger.GetInstance()
	logger.Println("Start update user ...")
	dbgoblog := connection.GetGoblogConn()

	var mess string
	var UserData []User

	sqlStat := fmt.Sprintf(`UPDATE users SET name = '%v', age = '%v', location = '%v' WHERE userid = '%v'`, usr.Name, usr.Age, usr.Location, id)
	conn, _ := dbgoblog.Begin()
	res, err := conn.Prepare(sqlStat)

	if config.FancyHandleError(err) {
		logger.Println("Failed update user, " + err.Error())
		logger.Println(sqlStat)
		conn.Rollback()
		response := ResponseUser{
			Message: "Failed update user, " + err.Error(),
			Data: nil,
		}
		return response
	}

	result, er := res.Exec()

	//check how many rows affected
	rowsAffected, f := result.RowsAffected()
	if f != nil {
		logger.Printf("Error while checking the affected rows  in update user with id : %v, error : %v", id, f)
	} else {
		logger.Printf("Total rows/records affected in update user with id : %v is : %v", id, rowsAffected)
	}

	defer res.Close()

	if config.FancyHandleError(er) {
		logger.Println("Failed update user, " + er.Error())
		conn.Rollback()
		response := ResponseUser{
			Message: "Failed update user, " + er.Error(),
			Data: nil,
		}
		return response
	}

	usr.ID = id
	conn.Commit()

	if rowsAffected > 0 {
		mess = fmt.Sprintf("Success, %v data was updated by userid : %v", rowsAffected, id)
	} else {
		mess = fmt.Sprintf("No data was updated by userid : %v", id)
	}

	response := ResponseUser{
		Message: mess,
		Data: UserData,
	}
	return response
}

func DeleteUser(id int64) ResponseUser {
	logger := customlogger.GetInstance()
	logger.Printf("Starting delete user with id : %v", id)
	var mess string
	var userData []User

	dbgoblog := connection.GetGoblogConn()
	sqlStat := `DELETE FROM users WHERE userid = ?`

	res, err := dbgoblog.Exec(sqlStat, id)

	if config.FancyHandleError(err) {
		logger.Printf("Failed in delete user with id : %v, error : %v", id, err)
		response := ResponseUser{
			Message:  fmt.Sprintf("Failed in delete user with id : %v, error : %v", id, err),
			Data: nil,
		}
		return response
	}

	rowsAffected, f := res.RowsAffected()

	if config.FancyHandleError(f) {
		logger.Printf("Error while checking the affected rows in delete with id : %v, error : %v", id, f)
		mess = fmt.Sprintf("Error while checking the affected rows in delete with id : %v, error : %v", id, f)
	} else if rowsAffected == 0 {
		logger.Printf("No data user was deleted by id : %v", id)
		mess = fmt.Sprintf("No data user was deleted by id : %v", id)
	} else {
		logger.Printf("Total rows/records affected in delete user with id : %v is : %v", id, rowsAffected)
		mess = fmt.Sprintf("Total rows/records affected in delete user with id : %v is : %v", id, rowsAffected)
	}

	response := ResponseUser{
		Message: mess,
		Data: userData,
	}
	return response
}