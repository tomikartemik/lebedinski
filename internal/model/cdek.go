package model

// Pvz представляет пункт выдачи заказов СДЭК
type Pvz struct {
	Code        string   `json:"code"`         // Код ПВЗ
	Name        string   `json:"name"`         // Наименование ПВЗ
	Address     string   `json:"address_full"` // Полный адрес (Уточнить поле в API)
	WorkTime    string   `json:"work_time"`    // Время работы
	Phones      []Phone  `json:"phones"`       // Телефоны
	Note        string   `json:"note"`         // Примечание
        Latitude    float64  `json:"latitude"`   // Широта (Уточнить поле в API)
        Longitude   float64  `json:"longitude"`  // Долгота (Уточнить поле в API)
	Type        string   `json:"type"`         // Тип ПВЗ (PVZ, POSTAMAT)
	// Добавьте другие поля при необходимости, сверившись с документацией API СДЭК
        // Возможно, поля address_full, latitude, longitude находятся внутри вложенного объекта location
        // Необходимо проверить актуальную документацию API СДЭК v2 /deliverypoints
}

// Phone представляет номер телефона ПВЗ
type Phone struct {
	Number string `json:"number"`
}

// Location может потребоваться, если координаты или адрес вложены
// type Location struct {
// 	AddressFull string  `json:"address_full"`
// 	Latitude    float64 `json:"latitude"`
// 	Longitude   float64 `json:"longitude"`
// }

// Структура ответа API СДЭК для списка ПВЗ (может потребоваться адаптация)
// Обычно API возвращает массив объектов Pvz напрямую.
// Если API возвращает объект с полем, содержащим массив,
// нужно будет создать соответствующую структуру ответа. 

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