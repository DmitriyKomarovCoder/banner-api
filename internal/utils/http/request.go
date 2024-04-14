package http

import (
	"net/http"
	"strconv"

	"github.com/DmitriyKomarovCoder/banner-api/internal/entity"
	"github.com/gorilla/mux"
)

const (
	tokenAdmin = "admin"
	tokenUser  = "user"
)

func GetValueFromUrl(value string, r *http.Request) (int, error) {
	idS := mux.Vars(r)[value]
	id, err := strconv.Atoi(idS)
	if err != nil {
		return 0, entity.ErrorsGetPath
	}

	return id, nil
}

func GetAuthToken(r *http.Request) bool {
	token := r.Header.Get("token")
	if token == tokenAdmin {
		return true
	}

	return false
}
