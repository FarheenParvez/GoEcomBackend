package main

import (
	"database/sql"
	"log"

	"github.com/FarheenParvez/GOECOMAPI/cmd/api"
	"github.com/FarheenParvez/GOECOMAPI/config"
	"github.com/FarheenParvez/GOECOMAPI/db"
	"github.com/go-sql-driver/mysql"
)

func main() {
	db, err := db.NewMySQLStorage(mysql.Config{
		User: config.Envs.DBUser,
		Passwd: config.Envs.DBPassword,
		Addr: config.Envs.DBAddress,
		Net: "tcp",
		DBName: config.Envs.DBName,
		AllowNativePasswords: true,
		ParseTime: true,
	})
	if err!= nil {
		log.Fatal(err)
	}
	
	initStorage(db)

	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err!= nil {
		log.Fatal(err)
	}

}

func initStorage(db *sql.DB) {
	err:= db.Ping()
	if err!= nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")

}