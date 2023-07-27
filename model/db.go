package model

import (
	"database/sql"
	"fmt"
	"log"
	"online_chess/util"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func GetDatabaseHandle() bool {
	db_user, user_ok := util.ReadEnvVariable("DB_USER")
	db_pwd, pwd_ok := util.ReadEnvVariable("DB_PWD")
	if !user_ok || !pwd_ok {
		fmt.Println("db env is empty!!!")
		return false
	}
	cfg := mysql.Config{
		User:   db_user,
		Passwd: db_pwd,
		Net:    "tcp",
		Addr:   "db:3306",
		DBName: "online_war_chess",
	}
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
		return false
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
		return false
	}
	return true
}

func GetDBInstance() *sql.DB {
	return db
}
