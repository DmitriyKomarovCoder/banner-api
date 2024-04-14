package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/DmitriyKomarovCoder/banner-api/pkg/logger"
	"github.com/gorilla/mux"
)

func PanicRecovery(logger *logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					logger.Error(string(debug.Stack()))
				}
			}()
			next.ServeHTTP(w, req)
		})
	}
}
