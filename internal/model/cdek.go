package model

// Pvz представляет пункт выдачи заказов СДЭК
type Pvz struct {
	Code     string   `json:"code"`     // Код ПВЗ
	Name     string   `json:"name"`     // Наименование ПВЗ
	Location Location `json:"location"` // Используем вложенную структуру Location
	WorkTime string   `json:"work_time"` // Время работы
	Phones   []Phone  `json:"phones"`   // Телефоны
	Note     string   `json:"note"`     // Примечание
	Type     string   `json:"type"`     // Тип ПВЗ (PVZ, POSTAMAT)
	// Убрали Latitude, Longitude, Address с верхнего уровня
	// Добавьте другие НЕ-location поля при необходимости
}

// Location содержит данные о местоположении ПВЗ
// Уточните поля по актуальной документации API СДЭК /deliverypoints!
type Location struct {
	CountryCode string  `json:"country_code"`
	RegionCode  int     `json:"region_code"`
	Region      string  `json:"region"`
	CityCode    int     `json:"city_code"`
	City        string  `json:"city"`
	FiasGUID    string  `json:"fias_guid,omitempty"` // Может отсутствовать
	PostalCode  string  `json:"postal_code"`
	Longitude   float64 `json:"longitude"` // Долгота
	Latitude    float64 `json:"latitude"`  // Широта
	Address     string  `json:"address"`      // Адрес внутри location
	AddressFull string  `json:"address_full"` // Полный адрес внутри location
}

// Phone представляет номер телефона ПВЗ
type Phone struct {
	Number string `json:"number"`
}

// CityInfo представляет информацию о городе из API /location/cities СДЭК
// Уточните поля по актуальной документации!
type CityInfo struct {
	Code        int     `json:"code"`         // Код населенного пункта СДЭК
	City        string  `json:"city"`         // Название населенного пункта
	FiasGUID    string  `json:"fias_guid"`    // ФИАС-код
	KladrCode   string  `json:"kladr_code"`   // КЛАДР-код
	CountryCode string  `json:"country_code"` // Код страны в формате ISO 3166-1 alpha-2
	Country     string  `json:"country"`      // Название страны
	Region      string  `json:"region"`       // Название региона
	RegionCode  int     `json:"region_code"`  // Код региона СДЭК
	Latitude    float64 `json:"latitude"`     // Широта
	Longitude   float64 `json:"longitude"`    // Долгота
	// ... могут быть и другие поля
} 