package http

import (
	"net/http"

	"encoding/json"

	"github.com/DmitriyKomarovCoder/banner-api/internal/entity"
	"github.com/DmitriyKomarovCoder/banner-api/pkg/logger"
)

func ErrorResponse(w http.ResponseWriter, code int, err error, message string, log logger.Logger) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	if code < http.StatusInternalServerError {
		log.Infof("invalid request: %v:", err)
	} else {
		log.Error(err.Error())
	}

	if code == http.StatusBadRequest || code == http.StatusInternalServerError {
		response := entity.ResponseError{ErrMsg: message}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Errorf("Error failed to marshal error message: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)

			if _, writeErr := w.Write([]byte("Can't encode error message into json, message: " + message)); writeErr != nil {
				log.Errorf("Error writing response: %s", writeErr.Error())
			}
		}
	}
}

func SuccessResponse[T any](w http.ResponseWriter, status int, response T) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
