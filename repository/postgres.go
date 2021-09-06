package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

func CreateDbConnection() (*gorm.DB, error) {
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	port := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	name := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, name, password, port)
	db, err = gorm.Open(dialect, dbURI)
	if err != nil {
		log.Println("error while creating db connection: ", err)
		return nil, err
	}
	err = db.DB().Ping()
	if err != nil {
		panic(err)
	}
	return db, nil
}
