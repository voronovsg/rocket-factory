package part

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/voronovsg/rocket-factory/inventory/internal/model"
	repoConv "github.com/voronovsg/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/voronovsg/rocket-factory/inventory/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, partUuid string) (model.Part, error) {
	var part repoModel.Part
	err := r.collection.FindOne(ctx, bson.M{partsFieldUUID: partUuid}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Part{}, model.ErrPartNotFound
		}
		return model.Part{}, err
	}

	return repoConv.PartToModel(part), nil
}
