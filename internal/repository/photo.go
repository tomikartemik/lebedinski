package repository

import (
	"gorm.io/gorm"
	"lebedinski/internal/model"
)

type PhotoRepository struct {
	db *gorm.DB
}

func NewPhotoRepository(db *gorm.DB) *PhotoRepository {
	return &PhotoRepository{db: db}
}

func (r *PhotoRepository) NewPhoto(photo model.Photo) error {
	return r.db.Create(&photo).Error
}
