package service

import (
	"lebedinski/internal/model"
	"lebedinski/internal/repository"
)

type PromoCodeService struct {
	repo repository.PromoCode
}

func NewPromoCodeService(repo repository.PromoCode) *PromoCodeService {
	return &PromoCodeService{repo: repo}
}

func (s *PromoCodeService) CreatePromoCode(promocode model.PromoCode) error {
	return s.repo.CreatePromoCode(promocode)
}

func (s *PromoCodeService) GetPromoCodeByCode(code string) (model.PromoCode, error) {
	promocode, err := s.repo.GetPromoCodeByCode(code)
	if promocode.NumberOfUses == 0 || err != nil {
		return model.PromoCode{}, err
	}
	return promocode, nil
}
