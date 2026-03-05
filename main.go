package main

import (
	"fmt"
	"net/http"
	"os"
	"web-blog/handlers"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()

	if _, err := os.Stat("articles"); os.IsNotExist(err) {
		os.Mkdir("articles", 0755)
	}
}

func main() {
	// Public pages
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/article/", handlers.ViewArticleHandler)

	// Auth
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	// Admin pages
	http.HandleFunc("/dashboard", handlers.DashboardArticleWithAuth())
	http.HandleFunc("/new", handlers.CreateArticleWithAuth())
	http.HandleFunc("/edit", handlers.UpdateArticleWithAuth())
	http.HandleFunc("/delete", handlers.DeleteArticleWithAuth())

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}