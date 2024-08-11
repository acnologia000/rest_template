package middleware

import (
	"net/http"
	"rest_template/shared/logger"
	"time"
)

func Chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func Logging(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Info("Request processed",
				map[string]interface{}{
					"method":   r.Method,
					"path":     r.URL.Path,
					"duration": time.Since(start),
				})
		})
	}
}
