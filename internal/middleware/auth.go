package middleware

import "net/http"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Токен авторизации не предоставлен", http.StatusUnauthorized)
			return
		}
		
		if len(token) < 10 {
			http.Error(w, "Неверный токен", http.StatusUnauthorized)
			return
		}
		
		next(w, r)
	}
}