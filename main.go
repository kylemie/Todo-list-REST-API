package main

import (
	"fmt"
	"rest-api/http"
	"rest-api/todo"
)

func main() {
	todoList := todo.NewList()
	httpHandlers := http.NewHTTPhandlers(todoList)
	httpserver := http.NewHTTPServer(httpHandlers)
	if err := httpserver.StartServer(); err != nil {
		fmt.Println("failed to start http server", err)
	}
}
