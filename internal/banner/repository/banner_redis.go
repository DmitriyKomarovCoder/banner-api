package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DmitriyKomarovCoder/banner-api/internal/entity"
	"github.com/go-redis/redis/v8"
)

type cache struct {
	db *redis.Client
}

func NewCache(db *redis.Client) *cache {
	return &cache{
		db: db,
	}
}

const (
	getCacheLayerMSG = "Get cache layer: %w"
	setCacheLayerMSG = "Set cache layer: %w"
	//==================
	TTL = 5 * time.Minute
)

func (r *cache) Set(tagID int, featureID int, content interface{}) error {
	key := fmt.Sprintf("%d:%d", tagID, featureID)

	jsonData, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	err = r.db.Set(context.Background(), key, jsonData, TTL).Err()
	if err != nil {
		return fmt.Errorf("failed to set data in cache: %v", err)
	}

	return nil
}

func (r *cache) Get(tagID int, featureID int) (interface{}, error) {
	key := fmt.Sprintf("%d:%d", tagID, featureID)

	content, err := r.db.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf(getCacheLayerMSG, entity.ErrorsNotFound)
		}
		return nil, fmt.Errorf(getCacheLayerMSG, err)
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(content), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %v", err)
	}

	return data, nil
}
