package part

import (
	"go.mongodb.org/mongo-driver/mongo"

	def "github.com/voronovsg/rocket-factory/inventory/internal/repository"
)

const (
	partsCollection               = "parts"
	partsFieldUUID                = "uuid"
	partsFieldName                = "name"
	partsFieldCategory            = "category"
	partsFieldTags                = "tags"
	partsFieldManufacturerCountry = "manufacturer.country"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *repository {
	return &repository{
		collection: db.Collection(partsCollection),
	}
}
