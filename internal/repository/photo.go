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

func (r *PhotoRepository) DeletePhoto(photoID int) error {
	var photo model.Photo
	if err := r.db.First(&photo, photoID).Error; err != nil {
		return err
	}
	return r.db.Delete(&photo).Error
}

func (r *PhotoRepository) GetPhotoByID(photoID int) (model.Photo, error) {
	var photo model.Photo
	if err := r.db.First(&photo, photoID).Error; err != nil {
		return photo, err
	}
	return photo, nil
}
