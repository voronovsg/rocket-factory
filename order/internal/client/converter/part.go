package converter

import (
	"log"
	"time"

	"github.com/voronovsg/rocket-factory/order/internal/model"
	inventoryV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/inventory/v1"
)

func PartListToModel(parts []*inventoryV1.Part) []model.Part {
	res := make([]model.Part, 0, len(parts))
	for _, part := range parts {
		res = append(res, PartToModel(part))
	}

	return res
}

func PartToModel(part *inventoryV1.Part) model.Part {
	var updatedAt *time.Time
	if part.UpdatedAt != nil {
		tmp := part.UpdatedAt.AsTime()
		updatedAt = &tmp
	}

	return model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      int32(part.Category),
		Dimensions:    DimensionsToModel(part.Dimensions),
		Manufacturer:  ManufacturerToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      PartMetadataToModel(part.Metadata),
		CreatedAt:     part.CreatedAt.AsTime(),
		UpdatedAt:     updatedAt,
	}
}

func DimensionsToModel(dimensions *inventoryV1.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ManufacturerToModel(manufacturer *inventoryV1.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func PartMetadataToModel(metadata map[string]*inventoryV1.Value) map[string]any {
	res := make(map[string]any, len(metadata))

	for key, value := range metadata {
		if value == nil || value.Kind == nil {
			continue
		}

		switch v := value.Kind.(type) {
		case *inventoryV1.Value_StringValue:
			res[key] = v.StringValue
		case *inventoryV1.Value_Int64Value:
			res[key] = v.Int64Value
		case *inventoryV1.Value_DoubleValue:
			res[key] = v.DoubleValue
		case *inventoryV1.Value_BoolValue:
			res[key] = v.BoolValue
		default:
			log.Printf("Unknown metadata kind for key %q: %T", key, value.Kind)
		}
	}

	return res
}

func PartsFilterToProto(filter model.PartsFilter) *inventoryV1.PartsFilter {
	categories := make([]inventoryV1.Category, 0, len(filter.Categories))
	for _, category := range filter.Categories {
		categories = append(categories, inventoryV1.Category(category))
	}

	return &inventoryV1.PartsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}
