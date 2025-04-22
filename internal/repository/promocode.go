package repository

import (
	"gorm.io/gorm"
	"lebedinski/internal/model"
)

type PromoCodeRepository struct {
	db *gorm.DB
}

func NewPromoCodeRepository(db *gorm.DB) *PromoCodeRepository {
	return &PromoCodeRepository{db: db}
}

func (r *PromoCodeRepository) CreatePromoCode(promocode model.PromoCode) error {
	return r.db.Create(&promocode).Error
}

func (r *PromoCodeRepository) GetPromoCodeByCode(code string) (model.PromoCode, error) {
	var promocode model.PromoCode

	if err := r.db.Where("code = ?", code).Find(&promocode).Error; err != nil {
		return promocode, err
	}
	return promocode, nil
}
