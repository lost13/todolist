package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	"io"
	"log"
	"net/http"
	"strconv"
)

var db *gorm.DB

func init() {
	_ = configor.Load(&Config, "config.yaml")
	fmt.Printf("config: %#v\n", Config)
}

func getdb() *gorm.DB {
	db, err := gorm.Open("postgres", "sslmode=disable host="+Config.Db.Host+" port="+Config.Db.Port+" user="+Config.Db.User+" dbname="+Config.Db.Name+" password="+Config.Db.Password+"")
	if err != nil {
		log.Panic(err)
	}
	return db
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	todo := &TodoModel{Description: description, Completed: false}
	println(todo.Completed)
	db.Create(&todo)
	result := db.Last(&todo)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result.Value)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := GetTodoByID(id)
	if err == false {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"update": false, "error": "ToDo not found"}`)
	} else {
		completed, _ := strconv.ParseBool(r.FormValue("completed"))
		todo := &TodoModel{}
		db.First(&todo, id)
		todo.Completed = completed
		println(todo.Completed)
		db.Save(&todo)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"update": true}`)
	}
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := GetTodoByID(id)
	if err == false {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"delete": false, "error": "ToDo not found"}`)
	} else {
		todo := &TodoModel{}
		db.First(&todo, id)
		db.Delete(&todo)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"delete": true}`)
	}
}

func GetTodos() interface{} {
	var todos []TodoModel
	TodoItems := db.Select("id").Find(&todos).Value
	return TodoItems
}

func GetTodoByID(Id int) bool {
	todo := &TodoModel{}
	result := db.First(&todo, Id)
	if result.Error != nil {
		return false
	}
	return true
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(GetTodos())
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.HandleFunc(route.Pattern, route.Handler).Methods(route.Method)
	}

	return router
}
