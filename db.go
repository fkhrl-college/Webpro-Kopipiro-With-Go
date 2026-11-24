package main

import (
	"database/sql"
	"errors"
	"log"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func connectDB() {
	dsn := "root:@tcp(127.0.0.1:3306)/db_school?charset=utf8mb4"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Connection failed: ", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Connection failed: ", err)
	}
}

// mysqlErrNo returns the MySQL error number inside err (0 if there is none).
// 1062 = duplicate value in a UNIQUE column
// 1452 = FOREIGN KEY check failed (the parent row does not exist)
func mysqlErrNo(err error) uint16 {
	if me, ok := errors.AsType[*mysql.MySQLError](err); ok {
		return me.Number
	}
	return 0
}
