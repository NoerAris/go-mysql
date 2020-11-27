package connection

import (
	"database/sql"
	"go-mysql/config"
	"time"
)

func GetGoblogConn () *sql.DB{
	dbgoblog := config.GetDbByPath("goblog").GetDb()
	dbgoblog.SetMaxIdleConns(0)
	dbgoblog.SetConnMaxLifetime(300 * time.Second)

	return dbgoblog
}
