package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/voronovsg/rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/inventory/v1"
)

func PartListToProto(parts []model.Part) []*inventoryV1.Part {
	res := make([]*inventoryV1.Part, 0, len(parts))
	for _, part := range parts {
		res = append(res, PartToProto(part))
	}

	return res
}

func PartToProto(part model.Part) *inventoryV1.Part {
	var updatedAt *timestamppb.Timestamp
	if part.UpdatedAt != nil {
		updatedAt = timestamppb.New(*part.UpdatedAt)
	}

	return &inventoryV1.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryToProto(part.Category),
		Dimensions:    DimensionsToProto(part.Dimensions),
		Manufacturer:  ManufacturerToProto(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataToProto(part.Metadata),
		CreatedAt:     timestamppb.New(part.CreatedAt),
		UpdatedAt:     updatedAt,
	}
}

func CategoryToProto(category int32) inventoryV1.Category {
	switch inventoryV1.Category(category) {
	case inventoryV1.Category_CATEGORY_UNSPECIFIED,
		inventoryV1.Category_CATEGORY_ENGINE,
		inventoryV1.Category_CATEGORY_FUEL,
		inventoryV1.Category_CATEGORY_PORTHOLE,
		inventoryV1.Category_CATEGORY_WING:
		return inventoryV1.Category(category)
	default:
		return inventoryV1.Category_CATEGORY_UNSPECIFIED
	}
}

func DimensionsToProto(dimensions model.Dimensions) *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ManufacturerToProto(manufacturer model.Manufacturer) *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func MetadataToProto(metadata map[string]any) map[string]*inventoryV1.Value {
	res := make(map[string]*inventoryV1.Value, len(metadata))
	for key, value := range metadata {
		switch v := value.(type) {
		case string:
			res[key] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_StringValue{StringValue: v},
			}
		case int64, int32, int:
			res[key] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_Int64Value{Int64Value: toInt64(v)},
			}
		case float32:
			res[key] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_DoubleValue{DoubleValue: float64(v)},
			}
		case float64:
			res[key] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_DoubleValue{DoubleValue: v},
			}
		case bool:
			res[key] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_BoolValue{BoolValue: v},
			}
		default:

		}
	}

	return res
}

func toInt64(value any) int64 {
	switch v := value.(type) {
	case int:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	default:
		return 0
	}
}

func PartsFilterToModel(filter *inventoryV1.PartsFilter) model.PartsFilter {
	if filter == nil {
		return model.PartsFilter{}
	}

	categories := make([]int32, 0, len(filter.Categories))
	for _, category := range filter.Categories {
		categories = append(categories, int32(category))
	}

	return model.PartsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}
