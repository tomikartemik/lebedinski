package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"lebedinski/internal/model"
	"lebedinski/internal/repository"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type CdekService struct {
	repoItem  repository.Item
	repoOrder repository.Order
}

func NewCdekService(itemRepo repository.Item, repoOrder repository.Order) *CdekService {
	return &CdekService{
		repoItem:  itemRepo,
		repoOrder: repoOrder,
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

func (s *CdekService) CreateCdekOrder(cartIDStr string) (string, error) {
	cartID, err := strconv.Atoi(cartIDStr)

	if err != nil {
		return "", err
	}

	order, err := s.repoOrder.GetOrderByCartID(cartID)
	if err != nil {
		return "", err
	}

	token, err := s.GetToken()
	if err != nil {
		return "", fmt.Errorf("failed to get CDEK token: %w", err)
	}

	shipmentPoint := os.Getenv("SHIPMENT_POINT")

	cdekReq := model.CdekOrderRequest{
		Number:     fmt.Sprint(order.CartID),
		TariffCode: 136,
		Recipient: model.CdekRecipient{
			Name: order.FullName,
			Phones: []model.CdekPhone{
				{Number: order.Phone},
			},
			Email: order.Email,
		},
		DeliveryPoint: order.PointCode,
		ShipmentPoint: shipmentPoint,

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

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusAccepted {
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

	order.CdekOrderUUID = cdekResp.Entity.UUID

	err = s.repoOrder.UpdateOrder(order)

	if err != nil {
		return "", err
	}

	return cdekResp.Entity.UUID, nil
}

func (s *CdekService) getCityCode(cityName string, countryCode string) (int, error) {
	token, err := s.GetToken()
	if err != nil {
		return 0, fmt.Errorf("failed to get CDEK token for city lookup: %w", err)
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetQueryParams(map[string]string{
			"country_codes": countryCode,
			"city":          cityName,
			"size":          "1",
		}).
		Get("https://api.cdek.ru/v2/location/cities")

	if err != nil {
		return 0, fmt.Errorf("failed to query CDEK cities API: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return 0, fmt.Errorf("CDEK cities API error: Status %s, Body: %s", resp.Status(), resp.String())
	}

	var cities []model.CityInfo
	if err := json.Unmarshal(resp.Body(), &cities); err != nil {
		return 0, fmt.Errorf("failed to unmarshal cities response: %w. Body: %s", err, resp.String())
	}

	if len(cities) == 0 {
		return 0, fmt.Errorf("city not found in CDEK database: %s, country: %s", cityName, countryCode)
	}

	log.Printf("Найден код города СДЭК: %d для %s (%s)", cities[0].Code, cityName, countryCode)
	return cities[0].Code, nil
}

func (s *CdekService) getRegionCode(regionName string, countryCode string) (int, error) {
	token, err := s.GetToken()
	if err != nil {
		return 0, fmt.Errorf("failed to get CDEK token for region lookup: %w", err)
	}

	client := resty.New()
	queryParams := map[string]string{
		"country_codes": countryCode,
		"region":        regionName,
		"size":          "5",
	}
	log.Printf("Запрос кода региона к API СДЭК /v2/location/regions с параметрами: %+v", queryParams)

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetQueryParams(queryParams).
		Get("https://api.cdek.ru/v2/location/regions")

	if err != nil {
		return 0, fmt.Errorf("failed to query CDEK regions API: %w", err)
	}

	bodyBytes := resp.Body()
	log.Printf("Ответ от API СДЭК /v2/location/regions (статус %d): %s", resp.StatusCode(), string(bodyBytes))

	if resp.StatusCode() != http.StatusOK {
		return 0, fmt.Errorf("CDEK regions API error: Status %d", resp.StatusCode())
	}

	var regions []model.RegionInfo
	if err := json.Unmarshal(bodyBytes, &regions); err != nil {
		return 0, fmt.Errorf("failed to unmarshal regions response: %w", err)
	}

	if len(regions) == 0 {
		return 0, fmt.Errorf("region not found in CDEK database (empty list returned): %s, country: %s", regionName, countryCode)
	}

	found := false
	var foundCode int
	targetRegionLower := strings.ToLower(regionName)

	log.Printf("Ищем регион '%s' в ответе API:", targetRegionLower)
	for i, reg := range regions {
		log.Printf("  Проверяем [%d]: Код=%d, Название='%s', Страна=%s", i, reg.RegionCode, reg.Region, reg.CountryCode)
		if strings.ToLower(reg.Region) == targetRegionLower {
			log.Printf("    Найдено точное совпадение! Используем код: %d", reg.RegionCode)
			foundCode = reg.RegionCode
			found = true
			break
		}
	}

	if !found {
		log.Printf("Точное совпадение для региона '%s' не найдено в ответе API.", regionName)
		return 0, fmt.Errorf("region not found in CDEK database (no exact match): %s, country: %s", regionName, countryCode)
	}

	log.Printf("Используем найденный код региона СДЭК: %d для %s (%s)", foundCode, regionName, countryCode)
	return foundCode, nil
}

func (s *CdekService) GetPvzList(params map[string]string) ([]model.Pvz, error) {
	cityCode := params["city_code"]
	countryCode := params["country_codes"]

	if cityCode == "" || countryCode == "" {
		return nil, errors.New("city code and country code are required in service params")
	}

	token, err := s.GetToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get CDEK token for PVZ list: %w", err)
	}

	pvzParams := map[string]string{
		"city_code": cityCode,
		"type":      "PVZ",
	}

	log.Printf("Запрос списка ПВЗ с параметрами для API СДЭК: %+v", pvzParams)

	client := resty.New()
	request := client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(pvzParams)

	resp, err := request.Get("https://api.cdek.ru/v2/deliverypoints")

	if err != nil {
		return nil, fmt.Errorf("failed to get PVZ list from CDEK API: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("CDEK deliverypoints API error: Status %s, Body: %s", resp.Status(), resp.String())
	}

	var pvzList []model.Pvz
	if err := json.Unmarshal(resp.Body(), &pvzList); err != nil {
		log.Printf("Ошибка разбора JSON ответа ПВЗ: %v. Тело ответа: %s", err, resp.String())
		return nil, fmt.Errorf("failed to unmarshal PVZ list response: %w", err)
	}

	return pvzList, nil
}
