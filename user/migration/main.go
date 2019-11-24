package main

import (
	"log"
	"msc/db"
	"msc/user"
)

func main() {
	db, cancelFunc, err := db.OpenDb()
	if err != nil {
		log.Fatal(err)
	}
	defer cancelFunc()
	db.AutoMigrate(&user.User{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&user.User{})
}
