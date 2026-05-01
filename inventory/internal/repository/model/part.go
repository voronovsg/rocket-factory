package model

import "time"

type Dimensions struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}

type Part struct {
	Uuid          string         `bson:"uuid"`
	Name          string         `bson:"name"`
	Description   string         `bson:"description"`
	Price         float64        `bson:"price"`
	StockQuantity int64          `bson:"stock_quantity"`
	Category      int32          `bson:"category"`
	Dimensions    Dimensions     `bson:"dimensions"`   // Dimensions представляет размеры детали
	Manufacturer  Manufacturer   `bson:"manufacturer"` // Manufacturer структура для хранения информации о производителе детали
	Tags          []string       `bson:"tags"`         // Tags теги для быстрого поиска
	Metadata      map[string]any `bson:"metadata"`     // Metadata гибкие метаданные
	CreatedAt     time.Time      `bson:"created_at"`
	UpdatedAt     *time.Time     `bson:"updated_at,omitempty"`
}
