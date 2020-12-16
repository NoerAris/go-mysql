package config

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
	// _ "gopkg.in/goracle.v2"
)

func init() {

}

type DBConf struct {
	username, password, host, port, schema, dbType, dbMigrateScript string
	db                                                              *sql.DB
}

func (conf *DBConf) getDbConnectionString() (string, error) {
	var connectionString string
	if conf.dbType == "postgres" {
		connectionString = fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", conf.username, conf.password, conf.host, conf.port, conf.schema)
	} else if conf.dbType == "mysql" {
		connectionString = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?multiStatements=true&parseTime=true", conf.username, conf.password, conf.host, conf.port, conf.schema)
	} else if conf.dbType == "oracle" {
		connectionString = fmt.Sprintf("oracle://%v:%v@%v:%v/%v", conf.username, conf.password, conf.host, conf.port, conf.schema)
	} else {
		log.Println(conf)
		return "", errors.New("unsupported db Type")
	}
	return connectionString, nil
}

func (conf *DBConf) loadDb() *sql.DB {
	if conf.db == nil {
		connectionString, err := conf.getDbConnectionString()
		if err != nil {
			log.Println(err.Error())
		}
		//log.Println("initialize connection with connection string " + connectionString)
		var newDb *sql.DB
		if conf.dbType == "oracle" {
			newDb, err = sql.Open("goracle", connectionString)
			newDb.SetMaxIdleConns(20)
			newDb.SetMaxOpenConns(100)
		} else {
			newDb, err = sql.Open(conf.dbType, connectionString)
			newDb.SetMaxIdleConns(20)
			newDb.SetMaxOpenConns(100)
		}
		if err != nil {
			log.Println(err.Error())
		}
		if conf.dbType == "oracle" {
			rows, err := newDb.Query("select 1 from dual")
			if err != nil {
				log.Println(err.Error())
				rows.Close()
			}
			defer rows.Close()
			ticker := time.NewTicker(time.Minute * 5)
			go func() {
				for _ = range ticker.C {
					rows, err := newDb.Query("select 1 from dual")
					if err != nil {
						log.Println("ERROR|KEEP ALIVE|", conf.dbType, err.Error())
						return
					}
					rows.Next()
					var i int
					rows.Scan(&i)
					log.Println("Keep Alive Connection ", i, conf.dbType)
					rows.Close()
				}
			}()

		} else {
			rows, err := newDb.Query("select 1")
			if err != nil {
				log.Println(err.Error())
				rows.Close()
			}
			ticker := time.NewTicker(time.Minute * 5)
			defer rows.Close()
			go func() {
				for _ = range ticker.C {
					rows, err := newDb.Query("select 1")
					if err != nil {
						log.Println("ERROR|KEEP ALIVE|", conf.dbType, err.Error())
						return
					}

					rows.Next()
					var i int
					rows.Scan(&i)
					log.Println("Keep Alive Connection ", i, conf.dbType)
					rows.Close()
				}
			}()
		}
		if err != nil {
			log.Println(err.Error())
		}
		log.Println("successfully connected to db", conf.dbType, conf.host)
		conf.db = newDb
	}
	return conf.db
}

func GetDefaultDb() *DBConf {
	db.username = GetValueFromConfig(env, "db."+INI_KEY_DB_USERNAME)
	db.password = GetValueFromConfig(env, "db."+INI_KEY_DB_PASSWORD)
	db.host = GetValueFromConfig(env, "db."+INI_KEY_DB_HOST)
	db.port = GetValueFromConfig(env, "db."+INI_KEY_DB_PORT)
	db.schema = GetValueFromConfig(env, "db."+INI_KEY_DB_SCHEMA)
	db.dbType = GetValueFromConfig(env, "db."+INI_KEY_DB_TYPE)
	return &db
}
func GetDbByPath(path string) *DBConf {
	if val, ok := dbs[path]; ok {
		//do something here
		return val
	} else {
		newdb := DBConf{}
		newdb.username = GetValueFromConfig(env, strings.Join([]string{"db", path, INI_KEY_DB_USERNAME}, "."))
		newdb.password = GetValueFromConfig(env, strings.Join([]string{"db", path, INI_KEY_DB_PASSWORD}, "."))
		newdb.host = GetValueFromConfig(env, strings.Join([]string{"db", path, INI_KEY_DB_HOST}, "."))
		newdb.port = GetValueFromConfig(env, strings.Join([]string{"db", path, INI_KEY_DB_PORT}, "."))
		newdb.schema = GetValueFromConfig(env, strings.Join([]string{"db", path, INI_KEY_DB_SCHEMA}, "."))
		newdb.dbType = GetValueFromConfig(env, strings.Join([]string{"db", path, INI_KEY_DB_TYPE}, "."))
		dbs[path] = &newdb
		return dbs[path]
	}

}

var db DBConf
var dbs = make(map[string]*DBConf)

func (db *DBConf) GetDb() *sql.DB {
	return db.loadDb()
}
