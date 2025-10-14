package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Serikkambarov/todo-scheduler/pkg/api"
	"github.com/Serikkambarov/todo-scheduler/pkg/db"
)

func main() {
	webDir := "./web"
	port := "7540"
	dbFile := "scheduler.db"

	// Инициализация базы данных
	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database initialized:", dbFile)

	// Инициализация API
	api.Init()

	// Настройка файлового сервера для фронтенда
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	fmt.Printf("Server started on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
