package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/voronovsg/rocket-factory/inventory/internal/config/env"
)

var appConfig *config

type config struct {
	InventoryGRPC InventoryGRPCConfig
	InventoryHTTP InventoryHTTPConfig
	Mongo         MongoConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	inventoryHTTPCfg, err := env.NewInventoryHTTPConfig()
	if err != nil {
		return err
	}

	inventoryGRPCCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	mongoCfg, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		InventoryGRPC: inventoryGRPCCfg,
		InventoryHTTP: inventoryHTTPCfg,
		Mongo:         mongoCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
