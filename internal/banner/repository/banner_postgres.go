package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DmitriyKomarovCoder/banner-api/internal/entity"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	checkTagsExistMSG    = "CheckIfTagsExist repository layer: %w"
	checkFeatureExistMSG = "CheckIfFeatureIdExist repository layer: %w"
	getBannerByIdMSG     = "GetBannerById repository layer: %w"
	getBannersMSG        = "GetBanners repository layer: %w"
	createBannerMSG      = "CreateBanner repository layer: %w"
	updateBannerMSG      = "UpdateBanner repository layer: %w"
	rmBannerMSG          = "DeleteBanner repository layer: %w"
	// =============================
	checkTags = `SELECT COUNT(*) 
				 FROM tags 
				 WHERE tag_id = $1;`

	checkFeature = `SELECT COUNT(*) 
					FROM features 
					WHERE feature_id = $1;`

	getBannerById = `SELECT banner_id, content, active, feature_id, created_at, update_at
				 FROM banners
				 WHERE banner_id = $1;`

	getTags = `SELECT tags.tag_id
			   FROM banners
			   JOIN banner_tags ON banners.banner_id = banner_tags.banner_id
			   JOIN tags ON banner_tags.tag_id = tags.tag_id
			   WHERE banners.banner_id = $1;`

	getBanner = `SELECT b.content, b.active
				 FROM banners b
				 JOIN features f ON b.feature_id = f.feature_id
				 JOIN banner_tags bt ON b.banner_id = bt.banner_id
				 JOIN tags t ON bt.tag_id = t.tag_id
				 WHERE b.feature_id = $1 AND bt.tag_id = $2
				 AND (b.active = true OR $3 = true)
				 ORDER BY b.update_at ASC
				 LIMIT 1;`

	getBanners = `
				 SELECT 
					 b.banner_id, 
					 ARRAY(SELECT bt.tag_id FROM banner_tags bt WHERE bt.banner_id = b.banner_id) AS tag_ids, 
					 b.feature_id, 
					 b.content, 
					 b.active, 
					 b.created_at, 
					 b.update_at 
				 FROM 
					 banners b 
				 WHERE 
					 (b.active = true OR $1 = true)
			 `

	createBannerSQL = `INSERT INTO banners (content, active, feature_id, created_at) 
						VALUES ($1, $2, $3, $4)
						RETURNING banner_id;`

	createBannerTagsSQL = `INSERT INTO banner_tags (banner_id, tag_id) VALUES ($1, $2);`

	rmBannerTagSQL = `DELETE FROM banner_tags WHERE banner_id = $1;`
	updBannerSQL   = `UPDATE banners SET content = $1, active = $2, feature_id = $3, update_at = $4 WHERE banner_id = $5;`
	rmBannerSQL    = `DELETE FROM banners WHERE banner_id = $1;`
)

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetBanner(tagId, featureId int, useLastRevision, isAdmin bool) (interface{}, bool, error) {
	var content interface{}

	active := false
	err := r.db.QueryRow(context.Background(), getBanner, featureId, tagId, isAdmin).Scan(&content, &active)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, entity.ErrorsNotFound
		}
		return nil, active, err
	}

	return content, active, nil
}

func (r *repository) GetBanners(tagId, featureId int, limit, offset int, isAdmin bool) ([]entity.Banner, error) {
	args := []interface{}{}
	query := getBanners

	args = append(args, isAdmin)
	count := 2
	if tagId != 0 {
		query += " AND EXISTS (SELECT 1 FROM banner_tags bt WHERE b.banner_id = bt.banner_id AND bt.tag_id = $" + fmt.Sprint(count) + ")"
		count++
		args = append(args, tagId)
	}

	if featureId != 0 {
		query += " AND b.feature_id = $" + fmt.Sprint(count)
		count++
		args = append(args, featureId)
	}

	query += " GROUP BY b.banner_id"
	if limit != 0 {
		query += " LIMIT $" + fmt.Sprint(count)
		count++
		args = append(args, limit)
	}

	if offset != 0 {
		query += " OFFSET $" + fmt.Sprint(count)
		count++
		args = append(args, offset)
	}

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrorsNotFound
		}
		return nil, err
	}
	defer rows.Close()

	banners := []entity.Banner{}
	for rows.Next() {
		var banner entity.Banner
		var tagIds []int
		if err := rows.Scan(&banner.BannerId, &tagIds, &banner.FeatureId, &banner.Content, &banner.IsActive, &banner.CreatedDate, &banner.UpdateDate); err != nil {
			return nil, fmt.Errorf(getBannersMSG, err)
		}

		banner.TagsId = tagIds

		banners = append(banners, banner)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(getBannersMSG, err)
	}

	return banners, nil
}

func (r *repository) CreateBanner(createBanner *entity.Banner) (int, error) {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return 0, fmt.Errorf(createBannerMSG, err)
	}
	defer tx.Rollback(context.Background())

	var bannerId int
	err = tx.QueryRow(context.Background(), createBannerSQL, createBanner.Content, createBanner.IsActive, createBanner.FeatureId, time.Now()).Scan(&bannerId)
	if err != nil {
		return 0, fmt.Errorf(createBannerMSG, err)
	}

	for _, tagId := range createBanner.TagsId {
		_, err = tx.Exec(context.Background(), createBannerTagsSQL, bannerId, tagId)
		if err != nil {
			return 0, fmt.Errorf(createBannerMSG, err)
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return 0, fmt.Errorf(createBannerMSG, err)
	}

	return bannerId, nil
}

func (r *repository) UpdateBanner(updBanner *entity.Banner) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return fmt.Errorf(updateBannerMSG, err)
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), rmBannerTagSQL, updBanner.BannerId)
	if err != nil {
		return fmt.Errorf(updateBannerMSG, err)
	}

	_, err = tx.Exec(context.Background(), updBannerSQL, updBanner.Content, updBanner.IsActive, updBanner.FeatureId, time.Now(), updBanner.BannerId)

	for _, tagId := range updBanner.TagsId {
		_, err = tx.Exec(context.Background(), createBannerTagsSQL, updBanner.BannerId, tagId)
		if err != nil {
			return fmt.Errorf(updateBannerMSG, err)
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf(updateBannerMSG, err)
	}

	return nil
}

func (r *repository) DeleteBanner(bannerId int) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return fmt.Errorf(rmBannerMSG, err)
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), rmBannerTagSQL, bannerId)
	if err != nil {
		return fmt.Errorf(rmBannerMSG, err)
	}

	_, err = tx.Exec(context.Background(), rmBannerSQL, bannerId)
	if err != nil {
		return fmt.Errorf(rmBannerMSG, err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf(rmBannerMSG, err)
	}

	return nil
}

func (r *repository) CheckIfTagsExist(tagIds []int) (bool, error) {
	if len(tagIds) == 0 {
		return false, fmt.Errorf(checkTagsExistMSG, entity.ErrorsNotFound)
	}

	var count, countRow int
	for _, id := range tagIds {
		err := r.db.QueryRow(context.Background(), checkTags, id).Scan(&countRow)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return false, fmt.Errorf(checkTagsExistMSG, entity.ErrorsNotFound)
			}
			return false, fmt.Errorf(checkTagsExistMSG, err)
		}
		count += countRow
	}

	return count == len(tagIds), nil
}

func (r *repository) CheckIfFeatureIdExist(featureId int) (bool, error) {
	var count int
	err := r.db.QueryRow(context.Background(), checkFeature, featureId).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf(checkFeatureExistMSG, entity.ErrorsNotFound)
		}
		return false, fmt.Errorf(checkFeatureExistMSG, err)
	}
	return count != 0, nil
}

func (r *repository) GetBannerById(bannerId int) (*entity.Banner, error) {
	var banner entity.Banner
	err := r.db.QueryRow(context.Background(), getBannerById, bannerId).Scan(
		&banner.BannerId,
		&banner.Content,
		&banner.IsActive,
		&banner.FeatureId,
		&banner.CreatedDate,
		&banner.UpdateDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &banner, fmt.Errorf(getBannerByIdMSG, entity.ErrorsNotFound)
		}
		return &banner, fmt.Errorf(getBannerByIdMSG, err)
	}

	rows, err := r.db.Query(context.Background(), getTags, bannerId)
	defer rows.Close()

	if err != nil {
		return &banner, fmt.Errorf(getBannerByIdMSG, err)
	}

	for rows.Next() {
		var tag int
		err := rows.Scan(&tag)
		if err != nil {
			return &banner, fmt.Errorf(getBannerByIdMSG, err)
		}
		banner.TagsId = append(banner.TagsId, tag)
	}

	if err := rows.Err(); err != nil {
		return &banner, fmt.Errorf(getBannerByIdMSG, err)
	}

	return &banner, nil
}
