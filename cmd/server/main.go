package main

import (
	"fmt"
	"net/http"

	"github.com/n1kno/practic8/internal/handler"
)

func main() {
	// Регистрируем маршруты
	http.HandleFunc("/", handler.FormHandler)
	http.HandleFunc("/greet", handler.GreetFormHandler)
	http.HandleFunc("/list", handler.ListHTMLHandler)
	http.HandleFunc("/api/greet", handler.GreetAPIHandler)
	http.HandleFunc("/api/list", handler.ListAPIHandler)

	fmt.Println("Сервис запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
