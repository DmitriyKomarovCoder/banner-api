package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/DmitriyKomarovCoder/banner-api/internal/banner"
	"github.com/DmitriyKomarovCoder/banner-api/internal/entity"
	"github.com/DmitriyKomarovCoder/banner-api/internal/entity/dto"
	util "github.com/DmitriyKomarovCoder/banner-api/internal/utils/http"
	"github.com/DmitriyKomarovCoder/banner-api/pkg/logger"
	"github.com/go-playground/validator/v10"
)

const (
	bannerIdPath = "id"
)

type Handler struct {
	usecase banner.Usecase
	log     logger.Logger
}

func NewHandler(usecase banner.Usecase, log logger.Logger) *Handler {
	return &Handler{
		usecase: usecase,
		log:     log,
	}
}

func (h *Handler) GetBanner(w http.ResponseWriter, r *http.Request) {
	isAdmin := util.GetAuthToken(r)
	tagIdS := r.URL.Query().Get("tag_id")
	featureIdS := r.URL.Query().Get("feature_id")
	lastRevisionS := r.URL.Query().Get("use_last_revision")

	if tagIdS == "" || featureIdS == "" {
		util.ErrorResponse(w, http.StatusBadRequest, nil, entity.MsgErrorQuery, h.log)
		return
	}

	tagId, err := strconv.Atoi(tagIdS)
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, err, entity.MsgErrorQuery, h.log)
		return
	}

	featureId, err := strconv.Atoi(featureIdS)
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, err, entity.MsgErrorQuery, h.log)
		return
	}

	var lastRevision bool
	if lastRevisionS != "" {
		lastRevision, err = strconv.ParseBool(lastRevisionS)
		if err != nil {
			util.ErrorResponse(w, http.StatusBadRequest, err, entity.MsgErrorQuery, h.log)
			return
		}
	} else {
		lastRevision = false
	}

	BannerContent, err := h.usecase.GetBanner(tagId, featureId, lastRevision, isAdmin)
	if err != nil {
		if errors.Is(err, entity.ErrorsNotFound) {
			h.log.Infof("invalid request: %v:", err)
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			h.log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	util.SuccessResponse(w, http.StatusOK, BannerContent)
}

func (h *Handler) GetBanners(w http.ResponseWriter, r *http.Request) {
	tagIdS := r.URL.Query().Get("tag_id")
	featureIdS := r.URL.Query().Get("feature_id")

	limitS := r.URL.Query().Get("limit")
	offsetS := r.URL.Query().Get("offset")

	var tagId int
	var err error

	if tagIdS != "" {
		tagId, err = strconv.Atoi(tagIdS)
		if err != nil {
			util.ErrorResponse(w, http.StatusBadRequest, err, entity.MsgErrorQuery, h.log)
			return
		}
	}

	var featureId int
	if featureIdS != "" {
		featureId, err = strconv.Atoi(featureIdS)
		if err != nil {
			util.ErrorResponse(w, http.StatusBadRequest, err, entity.MsgErrorQuery, h.log)
			return
		}
	}

	var limit, offset = 100, 0

	if limitS != "" {
		limit, err = strconv.Atoi(limitS)
		if err != nil {
			util.ErrorResponse(w, http.StatusBadRequest, err, entity.MsgErrorQuery, h.log)
			return
		}
	}

	if offsetS != "" {
		offset, err = strconv.Atoi(offsetS)
		if err != nil {
			util.ErrorResponse(w, http.StatusBadRequest, err, entity.MsgErrorQuery, h.log)
			return
		}
	}

	isAdmin := util.GetAuthToken(r)
	arrayBanner, err := h.usecase.GetBanners(tagId, featureId, limit, offset, isAdmin)
	if err != nil {
		if errors.Is(err, entity.ErrorsNotFound) {
			h.log.Infof("invalid request: %v:", err)
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			h.log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	bannersDTO := dto.BannerToArrayResponseDTO(arrayBanner)

	util.SuccessResponse(w, http.StatusOK, bannersDTO)
}

func (h *Handler) CreateBanners(w http.ResponseWriter, r *http.Request) {
	var BannerDTO dto.BannerCreateRequestDTO
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&BannerDTO); err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, err, entity.MsgErrorBody, h.log)
		return
	}

	validate := validator.New()
	if err := validate.Struct(BannerDTO); err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, err, entity.MsgErrorBody, h.log)
		return
	}

	banner := dto.BannerCreateDToToBanner(BannerDTO)
	bannerId, err := h.usecase.CreateBanner(&banner)
	if err != nil {
		if errors.Is(err, entity.ErrorsNotFound) {
			h.log.Infof("invalid request: %v:", err)
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			h.log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	util.SuccessResponse(w, http.StatusCreated, bannerId)
}

func (h *Handler) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	id, err := util.GetValueFromUrl(bannerIdPath, r)
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, err, entity.MsgErrorPath, h.log)
		return
	}

	var BannerDTO dto.BannerUpdateRequestDTO
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&BannerDTO); err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, err, entity.MsgErrorBody, h.log)
		return
	}
	banner := dto.BannerUpdateDToToBanner(BannerDTO, id)

	err = h.usecase.UpdateBanner(&banner)
	if err != nil {
		if errors.Is(err, entity.ErrorsNotFound) {
			h.log.Infof("invalid request: %v:", err)
			w.WriteHeader(http.StatusNotFound)
			return
		} else if errors.Is(err, entity.ErrorsNotBody) {
			h.log.Infof("invalid request: %v:", err)
			w.WriteHeader(http.StatusBadRequest)
		} else {
			h.log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	id, err := util.GetValueFromUrl(bannerIdPath, r)
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, err, entity.MsgErrorPath, h.log)
		return
	}

	err = h.usecase.DeleteBanner(id)
	if err != nil {
		if errors.Is(err, entity.ErrorsNotFound) {
			h.log.Infof("invalid request: %v:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			h.log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
