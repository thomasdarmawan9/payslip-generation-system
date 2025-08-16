package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	DEV         string = "dev"
	PROD        string = "prod"
	DEV_TEST    string = "dev_test"
	EnvFile            = "env/env.yml"
	EnvDevFile         = "env/env_dev.yml"
	EnvProdFile        = "env/env_prod.yml"
	ServiceCode        = "01"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[n.Int64()]
	}
	return string(result), nil
}
