package mysql

import (
	"fmt"

	"gorm.io/gorm"
)

var DB *gorm.DB

func DatabaseInit() {
	var err error

	// if we using mysql driver, reference : https://gorm.io/docs/connecting_to_the_database.html

	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	// DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to Database")
}
