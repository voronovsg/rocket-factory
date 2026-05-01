package part

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/voronovsg/rocket-factory/inventory/internal/model"
	repoConv "github.com/voronovsg/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/voronovsg/rocket-factory/inventory/internal/repository/model"
)

func (r *repository) List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	mongoFilter := bson.M{}

	if len(filter.Uuids) > 0 {
		mongoFilter[partsFieldUUID] = bson.M{"$in": filter.Uuids}
	}
	if len(filter.Names) > 0 {
		mongoFilter[partsFieldName] = bson.M{"$in": filter.Names}
	}
	if len(filter.Categories) > 0 {
		mongoFilter[partsFieldCategory] = bson.M{"$in": filter.Categories}
	}
	if len(filter.ManufacturerCountries) > 0 {
		mongoFilter[partsFieldManufacturerCountry] = bson.M{"$in": filter.ManufacturerCountries}
	}
	if len(filter.Tags) > 0 {
		mongoFilter[partsFieldTags] = bson.M{"$in": filter.Tags}
	}

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			log.Printf("failed to close cursor: %v\n", cerr)
		}
	}()
	var parts []repoModel.Part
	err = cursor.All(ctx, &parts)
	if err != nil {
		return nil, err
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return repoConv.ListPartsToModel(parts), nil
}
