package main

import (
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	"net/http"
	"os"
	"web-blog/middleware"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
	Success string `json:"success"`
}

func SetAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		MaxAge:   24 * 60 * 60,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
}

func hashPassword(pass string) string {
	hash := sha256.Sum256([]byte(pass))
	return hex.EncodeToString(hash[:])
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, nil)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminUsername == "" {
		adminUsername = "admin"
	}

	if adminPassword == "" {
		adminPassword = "123"
	}

	if username != adminUsername || hashPassword(password) != hashPassword(adminPassword) {
		w.WriteHeader(http.StatusUnauthorized)
		tmpl := template.Must(template.ParseFiles("templates/login_error.html"))
		tmpl.Execute(w, nil)
		return
	}

	token, err := middleware.GenerateJWT(username)
	if err != nil {
		http.Error(w, `{"error":"failed generate token"}`, http.StatusInternalServerError)
		return
	}

	SetAuthCookie(w,token)
	http.Redirect(w,r, "/dashboard", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
