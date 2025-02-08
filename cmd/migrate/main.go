package main

import (
	"log"
	"os"

	"github.com/FarheenParvez/GOECOMAPI/config"
	"github.com/FarheenParvez/GOECOMAPI/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	mysqlCfg "github.com/go-sql-driver/mysql"
)


func main() {
	db, err := db.NewMySQLStorage(mysqlCfg.Config{
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
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err!= nil {
		log.Fatal(err)
	}

	m , err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql", 
		driver, 
	)
	if err!= nil {
		log.Fatal(err)
	}


	cmd := os.Args[(len(os.Args) - 1)]
	if cmd == "up" {
		if err := m.Up(); err!= nil && err!= migrate.ErrNoChange {
			log.Fatal(err)
		}
	} else if cmd == "down" {
		if err := m.Down(); err!= nil && err!= migrate.ErrNoChange {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Invalid command")
	}

}