package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"lebedinski/internal/model"
	"lebedinski/internal/repository"
	"net/http"
	"os"
)

type CdekService struct {
	itemRepo repository.Item
}

func NewCdekService(itemRepo repository.Item) *CdekService {
	return &CdekService{
		itemRepo: itemRepo,
	}
}

func (s *CdekService) GetToken() (string, error) {
	account := os.Getenv("ACCOUNT_TOKEN")
	secure := os.Getenv("SECURE_TOKEN")
	client := resty.New()
	resp, err := client.R().
		SetFormData(map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     account,
			"client_secret": secure,
		}).
		Post("https://api.cdek.ru/v2/oauth/token")

	if err != nil {
		return "", err
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(resp.Body(), &tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

func (s *CdekService) GetCityCode(cityName string) (string, error) {
	token, err := s.GetToken()
	if err != nil {
		return "", err
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetQueryParams(map[string]string{
			"city":          cityName,
			"country_codes": "RU",
		}).
		Get("https://api.cdek.ru/v2/location/cities")

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		var errorResp struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(resp.Body(), &errorResp); err != nil {
			return "", fmt.Errorf("API error: %s", resp.String())
		}
		return "", fmt.Errorf("API error: %s", errorResp.Message)
	}

	var cities []struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(resp.Body(), &cities); err != nil {
		return "", err
	}

	if len(cities) == 0 {
		return "", errors.New("город не найден")
	}

	return cities[0].Code, nil
}

func (s *CdekService) CreateCdekOrder(order model.Order) (string, error) {
	token, err := s.GetToken()
	cityCode, err := s.GetCityCode(order.City)
	if err != nil {
		return "", err
	}

	cdekReq := &model.CdekOrderRequest{
		Number:     fmt.Sprint(order.CartID),
		TariffCode: 136,
		Recipient: struct {
			Name  string `json:"name"`
			Phone string `json:"phone"`
			Email string `json:"email"`
		}{
			Name:  order.FullName,
			Phone: order.Phone,
			Email: order.Email,
		},
		ToLocation: struct {
			Code    string `json:"code"`
			Address string `json:"address"`
			City    string `json:"city"`
			Country string `json:"country"`
		}{
			Code:    cityCode,
			Address: order.Address,
			City:    order.City,
			Country: "RU",
		},
		DeliveryCost: order.DeliveryCost,
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetBody(cdekReq).
		Post("https://api.cdek.ru/v2/orders")

	if err != nil {
		return "", err
	}

	var cdekResp struct {
		Entity struct {
			UUID string `json:"uuid"`
		} `json:"entity"`
	}
	if err := json.Unmarshal(resp.Body(), &cdekResp); err != nil {
		return "", err
	}

	return cdekResp.Entity.UUID, nil
}
