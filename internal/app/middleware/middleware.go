package middleware

import (
	"context"
	users "github.com/dantedoyl/car-life-api/internal/app/users"
	"net/http"
)

type Middleware struct {
	userUcase    users.IUsersUsecase
}

func NewMiddleware(userUcase users.IUsersUsecase) *Middleware {
	return &Middleware{
		userUcase: userUcase,
	}
}

func CorsControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:10888")

		switch req.Header.Get("Origin") {
		case "http://localhost:10888":
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:10888")

		case "http://127.0.0.1:10888":
			w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:10888")

		case "http://localhost:3000":
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

		case "http://127.0.0.1:3000":
			w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:3000")

		case "http://89.208.199.170":
			w.Header().Set("Access-Control-Allow-Origin", "http://89.208.199.170")
		}

		w.Header().Set("Access-Control-Allow-Headers", "content-type")
		w.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if req.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, req)
	})
}

func (m *Middleware) CheckAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		session, errE := m.userUcase.CheckSession(cookie.Value)
		if errE != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
