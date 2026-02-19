package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string
	Role     string
	jwt.RegisteredClaims
}

func JWTMIddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"malformed auth header"}`, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer")
		if tokenString == authHeader {
			http.Error(w, `{"error"}:"missing auth header"`, http.StatusUnauthorized)
			return
		}

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

func GenerateJWT(username string) (string, error) {
	claims := Claims{
		Username: username,
		Role:     "Admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func CookieAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("auth_token")
        if err != nil || cookie.Value == "" {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }

        claims, err := validateJWT(cookie.Value)
        if err != nil || claims == nil {
            clearAuthCookie(w)
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }

        ctx := context.WithValue(r.Context(), "user", claims)
        r = r.WithContext(ctx)
        next.ServeHTTP(w, r)
    }
}


func validateJWT(tokenString string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
    if err != nil || !token.Valid {
        return nil, err
    }
    return claims, nil
}

func clearAuthCookie(w http.ResponseWriter) {
    http.SetCookie(w, &http.Cookie{
        Name:     "auth_token",
        Value:    "",
        Path:     "/",
        MaxAge:   -1,
        HttpOnly: true,
    })
}

func GetUserFromContext(r *http.Request) *Claims {
    claims, ok := r.Context().Value("user").(*Claims)
    if !ok {
        return nil
    }
    return claims
}
