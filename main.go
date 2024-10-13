package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/scorcism/go-auth/cmd/api"
	"github.com/scorcism/go-auth/config"
	"github.com/scorcism/go-auth/db"
)

func main() {
	fmt.Println("Starting server")

	// setup db storage
	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal(err)
	}

	// Init connection
	initStorage(db)

	server := api.NewAPIServer(":8080", db)

	if error := server.Run(); error != nil {
		log.Fatal(error)
	}

}

func initStorage(db *sql.DB) {
	// connect to db
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: success connect")
}
