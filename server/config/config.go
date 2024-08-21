package config

import (
	"fmt"
	"os"
)

type Config struct {
	Address      string
	PostgresConn string
}

func ErrMissingEnvParam(param string) error {
	return fmt.Errorf("Missing required env param: %v", param)
}

func Get() (*Config, error) {
	addr := os.Getenv("ADDRESS")
	if addr == "" {
		return nil, ErrMissingEnvParam("ADDRESS")
	}

	pgConn := os.Getenv("POSTGRES_CONN")
	if pgConn == "" {
		return nil, ErrMissingEnvParam("POSTGRES_CONN")
	}

	return &Config{
		Address:      addr,
		PostgresConn: pgConn,
	}, nil

}
