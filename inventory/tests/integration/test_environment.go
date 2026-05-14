//go:build integration
// +build integration

package integration

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/bson"
)

// InsertTestPart — вставляет тестовую деталь в коллекцию Mongo и возвращает её UUID
func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUUID := gofakeit.UUID()

	partDoc := bson.M{
		"uuid":           partUUID,
		"name":           gofakeit.Word(),
		"description":    gofakeit.Sentence(5),
		"price":          gofakeit.Price(10, 1000),
		"stock_quantity": gofakeit.Number(1, 1000),
		"category":       1, // CATEGORY_ENGINE
		"manufacturer": bson.M{
			"name":    gofakeit.Company(),
			"country": gofakeit.Country(),
			"website": gofakeit.URL(),
		},
	}

	_, err := env.Mongo.Client().Database(env.Mongo.Config().Database).Collection(partsCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}

// ClearPartsCollection — удаляет все записи из коллекции parts
func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	_, err := env.Mongo.Client().Database(env.Mongo.Config().Database).Collection(partsCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
