package model

type Pvz struct {
	Code     string   `json:"code"`
	Name     string   `json:"name"`
	Location Location `json:"location"`
	Note     string   `json:"note"`
	Type     string   `json:"type"`
}

type Location struct {
	CountryCode string  `json:"country_code"`
	RegionCode  int     `json:"region_code"`
	Region      string  `json:"region"`
	CityCode    int     `json:"city_code"`
	City        string  `json:"city"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
	Address     string  `json:"address"`
}

type CityInfo struct {
	Code        int     `json:"code"`
	City        string  `json:"city"`
	CountryCode string  `json:"country_code"`
	Country     string  `json:"country"`
	Region      string  `json:"region"`
	RegionCode  int     `json:"region_code"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type RegionInfo struct {
	RegionCode  int    `json:"region_code"`
	Region      string `json:"region"`
	CountryCode string `json:"country_code"`
	Country     string `json:"country"`
}
