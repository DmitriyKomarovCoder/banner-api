package entity

import "time"

type Banner struct {
	BannerId    int
	TagsId      []int
	FeatureId   int
	Content     map[string]interface{}
	IsActive    *bool
	CreatedDate time.Time
	UpdateDate  time.Time
}
