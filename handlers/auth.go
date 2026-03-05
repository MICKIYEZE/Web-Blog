package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	"net/http"
	"os"
	"web-blog/handlers/middleware"
)

func hashPassword(pass string) string {
	h := sha256.Sum256([]byte(pass))
	return hex.EncodeToString(h[:])
}

func SetAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        tmpl := template.Must(template.ParseFiles("templates/login.html"))
        tmpl.Execute(w, nil)
        return
    }

    username := r.FormValue("username")
    password := r.FormValue("password")

    adminUser := os.Getenv("ADMIN_USERNAME")
    adminPass := os.Getenv("ADMIN_PASSWORD")

    if username != adminUser || password != adminPass {
        tmpl := template.Must(template.ParseFiles("templates/login_error.html"))
        tmpl.Execute(w, nil)
        return
    }

    token, _ := middleware.GenerateJWT(username)
    SetAuthCookie(w, token)

    http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	middleware.ClearAuthCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}