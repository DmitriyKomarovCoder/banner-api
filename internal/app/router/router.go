package router

import (
	banner "github.com/DmitriyKomarovCoder/banner-api/internal/banner/delivery/http"
	"github.com/DmitriyKomarovCoder/banner-api/pkg/logger"
	"github.com/DmitriyKomarovCoder/banner-api/pkg/middleware"
	"github.com/gorilla/mux"
)

func NewRouter(hBanner *banner.Handler, logger *logger.Logger) *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.PanicRecovery(logger))

	bannerRouter := r.PathPrefix("/api/v1").Subrouter()
	bannerRouter.Use(middleware.Auth)
	{
		bannerRouter.HandleFunc("/user_banner", hBanner.GetBanner).Methods("GET")
		bannerRouter.HandleFunc("/banner", hBanner.GetBanners).Methods("GET")
		bannerRouter.HandleFunc("/banner", hBanner.CreateBanners).Methods("POST")
		bannerRouter.HandleFunc("/banner/{id}", hBanner.UpdateBanner).Methods("PATCH")
		bannerRouter.HandleFunc("/banner/{id}", hBanner.DeleteBanner).Methods("DELETE")
	}

	return r
}
