package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"web-blog/handlers/middleware"
	"web-blog/model"
)

func parseTemplate(name string) *template.Template {
	return template.Must(template.ParseFiles("templates/" + name))
}

func getArticles() []model.Article {
	files, _ := os.ReadDir("articles")
	var list []model.Article

	for _, f := range files {
		if filepath.Ext(f.Name()) != ".json" {
			continue
		}

		data, _ := os.ReadFile("articles/" + f.Name())
		var a model.Article
		json.Unmarshal(data, &a)
		list = append(list, a)
	}
	return list
}

func getArticleByID(id int) *model.Article {
	path := fmt.Sprintf("articles/article%d.json", id)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var a model.Article
	json.Unmarshal(data, &a)
	return &a
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := parseTemplate("home.html")
	tmpl.Execute(w, getArticles())
}

func ViewArticleHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/article/")
	id, _ := strconv.Atoi(idStr)

	article := getArticleByID(id)
	if article == nil {
		http.NotFound(w, r)
		return
	}

	tmpl := parseTemplate("article.html")
	tmpl.Execute(w, article)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := parseTemplate("dashboard.html")
	tmpl.Execute(w, getArticles())
}

func createArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		parseTemplate("newArticle.html").Execute(w, nil)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	date := r.FormValue("date")

	if title == "" || content == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	if _, err := time.Parse("2006-01-02", date); err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	// Generate ID
	maxID := 0
	for _, a := range getArticles() {
		if a.ID > maxID {
			maxID = a.ID
		}
	}

	article := model.Article{
		ID:        maxID + 1,
		Title:     title,
		Content:   content,
		Published: date,
		Author:    user.Username,
	}

	file, _ := os.Create(fmt.Sprintf("articles/article%d.json", article.ID))
	json.NewEncoder(file).Encode(article)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func updateArticleHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	if r.Method == http.MethodGet {
		article := getArticleByID(id)
		parseTemplate("updateArticle.html").Execute(w, article)
		return
	}

	// POST update
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	date := r.FormValue("date")

	article := model.Article{
		ID:        id,
		Title:     title,
		Content:   content,
		Published: date,
		Author:    user.Username,
	}

	file, _ := os.Create(fmt.Sprintf("articles/article%d.json", id))
	json.NewEncoder(file).Encode(article)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func deleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	os.Remove(fmt.Sprintf("articles/article%d.json", id))
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Auth wrappers
func DashboardArticleWithAuth() http.HandlerFunc {
	return middleware.CookieAuthMiddleware(dashboardHandler)
}
func CreateArticleWithAuth() http.HandlerFunc {
	return middleware.CookieAuthMiddleware(createArticleHandler)
}
func UpdateArticleWithAuth() http.HandlerFunc {
	return middleware.CookieAuthMiddleware(updateArticleHandler)
}
func DeleteArticleWithAuth() http.HandlerFunc {
	return middleware.CookieAuthMiddleware(deleteArticleHandler)
}