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
	if err != nil {
		return model.PromoCode{}, err
	}
	if promocode.NumberOfUses == 0 {
		return model.PromoCode{}, err
	}
	return promocode, nil
}

func (s *PromoCodeService) GetAllPromoCodes() ([]model.PromoCode, error) {
	return s.repo.GetAllPromocodes()
}

func (s *PromoCodeService) DeletePromoCodeByCode(code string) error {
	return s.repo.DeletePromoCodeByCode(code)
}

func (s *PromoCodeService) UpdatePromoCode(promocode model.PromoCode) error {
	return s.repo.UpdatePromoCode(promocode)
}
