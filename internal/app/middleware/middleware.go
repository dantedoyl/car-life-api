package middleware

import "net/http"

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
