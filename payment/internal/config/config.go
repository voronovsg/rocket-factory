package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/voronovsg/rocket-factory/payment/internal/config/env"
)

var appConfig *config

type config struct {
	PaymentGRPC PaymentGRPCConfig
	PaymentHTTP PaymentHTTPConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	paymentGRPCCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	paymentHTTPCfg, err := env.NewPaymentHTTPConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		PaymentGRPC: paymentGRPCCfg,
		PaymentHTTP: paymentHTTPCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
