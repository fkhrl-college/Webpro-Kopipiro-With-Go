package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_school?charset=utf8mb4")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", "admin", string(hash), "admin")

	if me, ok := err.(*mysql.MySQLError); ok && me.Number == 1062 {
		fmt.Println("The admin account already exists. Nothing to do.")
		return
	} else if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Admin account created. Username: admin, Password: admin123. Delete this folder now!")
}
