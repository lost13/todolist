package main

import "net/http"

type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

type TodoModel struct {
	Id          int `gorm:"primary_key"`
	Description string
	Completed   bool
}

type Routes []Route

var routes = Routes{
	Route{
		Name:    "Index",
		Method:  "GET",
		Pattern: "/",
		Handler: TodoIndex,
	},
	Route{
		Name:    "TodoAdd",
		Method:  "GET",
		Pattern: "/todoadd",
		Handler: CreateTodo,
	},
	Route{
		Name:    "TodoUpdate",
		Method:  "GET",
		Pattern: "/todoup/{id}",
		Handler: UpdateTodo,
	},
	Route{
		Name:    "TodoDelete",
		Method:  "GET",
		Pattern: "/tododel/{id}",
		Handler: DeleteTodo,
	},
}

var Config = struct {
	APPName string `default:"todolist"`

	Db struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
		Host     string `default:"localhost"`
		Port     string `default:"5432"`
		Drop     bool   `default:"false"`
	}

}{}
