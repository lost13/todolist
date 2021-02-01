package main

import (
	_ "github.com/jackc/pgx"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
)

func main() {
	db = getdb()
	defer db.Close()

	if Config.Db.Drop == true {
		db.Debug().DropTableIfExists(&TodoModel{})
		db.Debug().AutoMigrate(&TodoModel{})
	}

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
