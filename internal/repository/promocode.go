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

func (r *PromoCodeRepository) GetAllPromocodes() ([]model.PromoCode, error) {
	var promocodes []model.PromoCode
	if err := r.db.Find(&promocodes).Error; err != nil {
		return promocodes, err
	}
	return promocodes, nil
}

func (r *PromoCodeRepository) DeletePromoCodeByCode(code string) error {
	return r.db.Delete(&model.PromoCode{}, "code = ?", code).Error
}

func (r *PromoCodeRepository) UpdatePromoCode(promocode model.PromoCode) error {
	return r.db.Save(&promocode).Error
}
