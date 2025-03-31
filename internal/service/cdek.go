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
	"strings"
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

	fmt.Println("DEBUG: Received CDEK Token:", tokenResp.AccessToken)
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
		fmt.Println(resp.Body())
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
	if err != nil {
		fmt.Println("ERROR: Failed to get CDEK token:", err)
		return "", fmt.Errorf("failed to get CDEK token: %w", err)
	}
	fmt.Println("DEBUG: Using CDEK Token for CreateOrder:", token)

	cdekReq := model.CdekOrderRequest{
		Number:     fmt.Sprint(order.CartID),
		TariffCode: 136,
		Recipient: struct {
			Name  string `json:"name"`
			Phone string `json:"phone"`
			Email string `json:"email,omitempty"`
		}{
			Name:  order.FullName,
			Phone: order.Phone,
			Email: order.Email,
		},
		DeliveryPoint: order.PointCode,
		Packages: []model.CdekPackage{
			{
				Number: fmt.Sprintf("%s-1", order.CartID),
				Weight: 1000,
				Length: 10,
				Width:  10,
				Height: 10,
				Items: []model.CdekPackageItem{
					{
						Name:    "Пример товара",
						WareKey: "ART-001",
						Payment: model.CdekPayment{
							Value: 0,
						},
						Cost:   1.0,
						Weight: 1000,
						Amount: 1,
					},
				},
			},
		},
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetBody(cdekReq).
		Post("https://api.cdek.ru/v2/orders")

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		errorMsg := fmt.Sprintf("CDEK API error: Status %s", resp.Status())
		var errorResp struct {
			Requests []struct {
				Errors []struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				} `json:"errors"`
				State string `json:"state"`
			} `json:"requests"`
		}
		if err := json.Unmarshal(resp.Body(), &errorResp); err == nil && len(errorResp.Requests) > 0 && len(errorResp.Requests[0].Errors) > 0 {
			var errorDetails []string
			for _, req := range errorResp.Requests {
				for _, e := range req.Errors {
					errorDetails = append(errorDetails, fmt.Sprintf("[%s] %s", e.Code, e.Message))
				}
			}
			if len(errorDetails) > 0 {
				errorMsg = fmt.Sprintf("%s. Details: %s", errorMsg, strings.Join(errorDetails, "; "))
			}
		} else {
			errorMsg = fmt.Sprintf("%s. Response Body: %s", errorMsg, resp.String())
		}
		return "", errors.New(errorMsg)
	}

	var cdekResp struct {
		Entity struct {
			UUID string `json:"uuid"`
		} `json:"entity"`
		Requests []struct {
			Errors []struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"errors"`
			State string `json:"state"`
		} `json:"requests"`
	}
	if err := json.Unmarshal(resp.Body(), &cdekResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal CDEK response: %w. Body: %s", err, resp.String())
	}

	if len(cdekResp.Requests) > 0 && len(cdekResp.Requests[0].Errors) > 0 {
		var errorDetails []string
		for _, req := range cdekResp.Requests {
			for _, e := range req.Errors {
				errorDetails = append(errorDetails, fmt.Sprintf("[%s] %s", e.Code, e.Message))
			}
		}
		return "", fmt.Errorf("CDEK returned success status but with errors: %s", strings.Join(errorDetails, "; "))
	}

	if cdekResp.Entity.UUID == "" {
		return "", fmt.Errorf("CDEK response successful, but UUID is empty. Body: %s", resp.String())
	}

	return cdekResp.Entity.UUID, nil
}
