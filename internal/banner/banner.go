package banner

import (
	"github.com/DmitriyKomarovCoder/banner-api/internal/entity"
)

// var _ Usecase = (*)(nil)

type Usecase interface {
	GetBanner(tagId, featureId int, useLastRevision, isAdmin bool) (interface{}, error)
	GetBanners(tagId, featureId int, limit, offset int, isAdmin bool) ([]entity.Banner, error)
	CreateBanner(createBanner *entity.Banner) (int, error)
	UpdateBanner(updBanner *entity.Banner) error
	DeleteBanner(bannerId int) error
}

// var _ Repository = (*test)(nil)
type Repository interface {
	GetBannerById(bannerId int) (*entity.Banner, error)
	GetBanner(tagId, featureId int, useLastRevision, isAdmin bool) (interface{}, bool, error)
	GetBanners(tagId, featureId int, limit, offset int, isAdmin bool) ([]entity.Banner, error)
	CreateBanner(createBanner *entity.Banner) (int, error)
	UpdateBanner(updBanner *entity.Banner) error
	DeleteBanner(bannerId int) error
	CheckIfTagsExist(tagIds []int) (bool, error)
	CheckIfFeatureIdExist(featureId int) (bool, error)
}

type Cashe interface {
	Set(tagID int, featureID int, content interface{}) error
	Get(tagID int, featureID int) (interface{}, error)
}
