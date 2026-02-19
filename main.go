package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Article struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Published string `json:"published"`
}

func createArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var a Article
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if a.Title == "" || len(a.Title) > 100 {
		http.Error(w, `{"error":"title required and <100 chars"}`, http.StatusBadRequest)
		return
	}

	if a.Content == "" {
		http.Error(w, `{"error":"content required"}`, http.StatusBadRequest)
		return
	}

	if _, err := time.Parse("2006-01-02", a.Published); err != nil {
		http.Error(w, `{"error":"invalid date format YYYY-MM-DD"}`, http.StatusBadRequest)
		return
	}

	files, _ := os.ReadDir("articles")
	a.ID = len(files) + 1

	filePath := fmt.Sprintf("articles/article%d.json", a.ID)
	file, _ := os.Create(filePath)
	defer file.Close()
	json.NewEncoder(file).Encode(a)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(a)
}

func getArticles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	files, _ := os.ReadDir("articles")
	var articles []Article

	for _, f := range files {
		if filepath.Ext(f.Name()) != ".json" {
			continue
		}
		data, _ := os.ReadFile(filepath.Join("articles", f.Name()))
		var art Article
		json.Unmarshal(data, &art)
		articles = append(articles, art)
	}

	json.NewEncoder(w).Encode(articles)
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/articles/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	var a Article
	err = json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	a.ID = id

	filePath := fmt.Sprintf("articles/article%d.json", a.ID)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	file, _ := os.Create(filePath)
	defer file.Close()
	json.NewEncoder(file).Encode(a)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(a)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/articles/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	filePath := fmt.Sprintf("articles/article%d.json", id)
	if err := os.Remove(filePath); err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	http.HandleFunc("/articles", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createArticle(w, r)
		} else if r.Method == http.MethodGet {
			getArticles(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/articles/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			updateArticle(w, r)
		} else if r.Method == http.MethodDelete {
			deleteArticle(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)
}
