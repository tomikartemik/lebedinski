package model

// Pvz представляет пункт выдачи заказов СДЭК
type Pvz struct {
	Code     string   `json:"code"`
	Name     string   `json:"name"`
	Location Location `json:"location"`
	WorkTime string   `json:"work_time"`
	Phones   []Phone  `json:"phones"`
	Note     string   `json:"note"`
	Type     string   `json:"type"`
}

type Location struct {
	CountryCode string  `json:"country_code"`
	RegionCode  int     `json:"region_code"`
	Region      string  `json:"region"`
	CityCode    int     `json:"city_code"`
	City        string  `json:"city"`
	FiasGUID    string  `json:"fias_guid,omitempty"`
	PostalCode  string  `json:"postal_code"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
	Address     string  `json:"address"`
	AddressFull string  `json:"address_full"`
}

type Phone struct {
	Number string `json:"number"`
}

type CityInfo struct {
	Code        int     `json:"code"`
	City        string  `json:"city"`
	FiasGUID    string  `json:"fias_guid"`
	KladrCode   string  `json:"kladr_code"`
	CountryCode string  `json:"country_code"`
	Country     string  `json:"country"`
	Region      string  `json:"region"`
	RegionCode  int     `json:"region_code"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	// ... могут быть и другие поля
}
