package dto

import (
	"time"

	"github.com/DmitriyKomarovCoder/banner-api/internal/entity"
)

type BannerContentResponseDTO struct {
	Content interface{}
}

func BannerToContentResponseDTO(banner entity.Banner) BannerContentResponseDTO {
	return BannerContentResponseDTO{
		Content: banner.Content,
	}
}

type BannerResponseDTO struct {
	BannerId    int                    `json:"banner_id"`
	TagsId      []int                  `json:"tag_ids"`
	FeatureId   int                    `json:"feature_id"`
	Content     map[string]interface{} `json:"content"`
	IsActive    *bool                  `json:"is_active"`
	CreatedDate time.Time              `json:"created_at"`
	UpdateDate  time.Time              `json:"updated_at"`
}

func BannerToArrayResponseDTO(banners []entity.Banner) []BannerResponseDTO {
	var bannersDTO []BannerResponseDTO
	for _, banner := range banners {
		bannersDTO = append(bannersDTO, BannerResponseDTO(banner))
	}
	return bannersDTO
}

type BannerCreateRequestDTO struct {
	TagIds    []int                  `json:"tag_ids" validate:"required"`
	FeatureId int                    `json:"feature_id" validate:"required"`
	Content   map[string]interface{} `json:"content" validate:"required"`
	IsActive  *bool                  `json:"is_active" validate:"required"`
}

func BannerCreateDToToBanner(bannerDTO BannerCreateRequestDTO) entity.Banner {
	return entity.Banner{
		TagsId:    bannerDTO.TagIds,
		FeatureId: bannerDTO.FeatureId,
		Content:   bannerDTO.Content,
		IsActive:  bannerDTO.IsActive,
	}
}

type BannerUpdateRequestDTO struct {
	TagIds    []int                  `json:"tag_ids"`
	FeatureId int                    `json:"feature_id"`
	Content   map[string]interface{} `json:"content"`
	IsActive  *bool                  `json:"is_active"`
}

func BannerUpdateDToToBanner(bannerDTO BannerUpdateRequestDTO, id int) entity.Banner {
	return entity.Banner{
		BannerId:  id,
		TagsId:    bannerDTO.TagIds,
		FeatureId: bannerDTO.FeatureId,
		Content:   bannerDTO.Content,
		IsActive:  bannerDTO.IsActive,
	}
}
