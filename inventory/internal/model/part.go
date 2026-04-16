package model

import "time"

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Part struct {
	Uuid          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      int32
	Dimensions    Dimensions     // Dimensions представляет размеры детали
	Manufacturer  Manufacturer   // Manufacturer cтруктура для хранения информации о производителе детали
	Tags          []string       // Tags теги для быстрого поиска
	Metadata      map[string]any // Metadata гибкие метаданные
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}

type PartsFilter struct {
	Uuids                 []string // Список UUID'ов деталей
	Names                 []string // Список наименований деталей
	Categories            []int32  // Список категорий деталей
	ManufacturerCountries []string // Страны производителей
	Tags                  []string // Теги для поиска
}
