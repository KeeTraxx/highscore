package main

import "github.com/jinzhu/gorm"
import _ "github.com/jinzhu/gorm/dialects/sqlite"
import "os"

var DB *gorm.DB

func initDB() {
	//db, err := gorm.Open("mysql", "highscore:highscore@tcp(localhost:3306)/highscore?charset=utf8&parseTime=True&loc=Local")
	err := os.MkdirAll("data", 0755)
	//os.RemoveAll("data/sqlite3.db")

	DB, err = gorm.Open("sqlite3", "data/sqlite3.db")

	if err != nil {
		panic(err)
	}

}
