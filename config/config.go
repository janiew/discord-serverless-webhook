package config

import "os"

func MustGetEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	if fallbackValue == "" {
		panic("Required env var not set: " + key)
	}

	return fallbackValue
}