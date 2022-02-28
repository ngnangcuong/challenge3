package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"fmt"
	"challenge3/models"
)

func GetDatabase() *gorm.DB {
	dbname := "Challenge3"
	db := "postgres"
	dbpassword := "Cuongnguyen2001"
	dburl := "postgres://postgres:" + dbpassword + "@localhost/" + dbname + "?sslmode=disable"

	connection, err := gorm.Open(db, dburl)

	if err != nil {
		log.Fatalln("Wrong database url")
	}

	sqldb := connection.DB()
	err = sqldb.Ping()

	if err != nil {
		log.Fatalln("Database is connectedt")
	}
	// log.Fatalln("Database is connected")
	fmt.Println("Database is connected")

	return connection
}

func InitMigration() {
	connection := GetDatabase()
	defer CloseDatabase(connection)

	connection.AutoMigrate(models.User{})
	connection.AutoMigrate(models.Post{})
	connection.AutoMigrate(models.Role{})

}

func CloseDatabase(connection *gorm.DB) {
	sqldb := connection.DB()
	sqldb.Close()
}