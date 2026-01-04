package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func GetString(key, fallback string) string {
	godotenv.Load()
	val, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return val
}

func GetInt(key string, fallback int) int {
	val, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return valAsInt
}
