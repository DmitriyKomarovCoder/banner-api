package usecase

import (
	"errors"
	"fmt"

	"github.com/DmitriyKomarovCoder/banner-api/internal/banner"
	"github.com/DmitriyKomarovCoder/banner-api/internal/entity"
)

type Usecase struct {
	bannerRepo  banner.Repository
	bannerCache banner.Cashe
}

func NewUsecase(br banner.Repository, bc banner.Cashe) *Usecase {
	return &Usecase{
		bannerRepo:  br,
		bannerCache: bc,
	}
}

const (
	getBannerMSG    = "GetBanner usecase layer: %w"
	getBannersMSG   = "GetBanners usecase layer: %w"
	createBannerMSG = "CreateBanner usecase layer: %w"
	updateBannerMSG = "UpdateBanner usecase layer: %w"
	deleteBannerMSG = "DeleteBanner usecase layer: %w"
)

func (u *Usecase) GetBanner(tagId, featureId int, useLastRevision, isAdmin bool) (interface{}, error) {
	var content interface{}
	var err error
	var active bool
	if useLastRevision {
		content, active, err = u.bannerRepo.GetBanner(tagId, featureId, useLastRevision, isAdmin)
		if err != nil {
			return nil, fmt.Errorf(getBannerMSG, err)
		}
		return content, nil
	}

	content, err = u.bannerCache.Get(tagId, featureId)
	if err != nil {
		if errors.Is(err, entity.ErrorsNotFound) {
			content, active, err = u.bannerRepo.GetBanner(tagId, featureId, useLastRevision, isAdmin)
			if err != nil {
				return nil, fmt.Errorf(getBannerMSG, err)
			}
			if active {
				err = u.bannerCache.Set(tagId, featureId, content)
			}

			if err != nil {
				return nil, fmt.Errorf(getBannerMSG, err)
			}
		} else {
			return nil, fmt.Errorf(getBannerMSG, err)
		}
	}

	return content, nil
}

func (u *Usecase) GetBanners(tagId, featureId int, limit, offset int, isAdmin bool) ([]entity.Banner, error) {
	banners, err := u.bannerRepo.GetBanners(tagId, featureId, limit, offset, isAdmin)
	if err != nil {
		return nil, fmt.Errorf(getBannersMSG, err)
	}
	return banners, nil
}

func (u *Usecase) CreateBanner(createBanner *entity.Banner) (int, error) {
	flag, err := u.bannerRepo.CheckIfTagsExist(createBanner.TagsId)
	if !flag || err != nil {
		return 0, fmt.Errorf(createBannerMSG, err)
	}

	flag, err = u.bannerRepo.CheckIfFeatureIdExist(createBanner.FeatureId)
	if !flag || err != nil {
		return 0, fmt.Errorf(createBannerMSG, err)
	}

	bannerId, err := u.bannerRepo.CreateBanner(createBanner)
	if err != nil {
		return 0, fmt.Errorf(createBannerMSG, err)
	}

	return bannerId, nil
}

func (u *Usecase) UpdateBanner(updBanner *entity.Banner) error {
	currentBanner, err := u.bannerRepo.GetBannerById(updBanner.BannerId)
	if err != nil {
		return fmt.Errorf(updateBannerMSG, err)
	}

	var flag bool
	if len(updBanner.TagsId) != 0 {
		flag, err = u.bannerRepo.CheckIfTagsExist(updBanner.TagsId)
		if !flag || err != nil {
			return fmt.Errorf(updateBannerMSG, err)
		}
	}

	if updBanner.FeatureId != 0 {
		flag, err = u.bannerRepo.CheckIfFeatureIdExist(updBanner.FeatureId)
		if !flag || err != nil {
			return fmt.Errorf(updateBannerMSG, err)
		}
	}

	if updBanner.Content == nil {
		updBanner.Content = currentBanner.Content
	}

	if len(updBanner.TagsId) == 0 {
		updBanner.TagsId = currentBanner.TagsId
	}

	if updBanner.FeatureId == 0 {
		updBanner.FeatureId = currentBanner.FeatureId
	}

	if updBanner.IsActive == nil {
		updBanner.IsActive = currentBanner.IsActive
	}

	if err := u.bannerRepo.UpdateBanner(updBanner); err != nil {
		return fmt.Errorf(updateBannerMSG, err)
	}

	return nil
}

func (r *Usecase) DeleteBanner(bannerId int) error {
	if err := r.bannerRepo.DeleteBanner(bannerId); err != nil {
		return fmt.Errorf(deleteBannerMSG, r.bannerRepo.DeleteBanner(bannerId))
	}
	return nil
}
