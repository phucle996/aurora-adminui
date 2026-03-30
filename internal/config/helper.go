package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getRequiredEnv(key string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		log.Fatalf("%s is required", key)
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	out, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("invalid %s: %v", key, err)
	}
	return out
}

func requirePairedEnv(left, right string) {
	leftValue := strings.TrimSpace(os.Getenv(left))
	rightValue := strings.TrimSpace(os.Getenv(right))
	if (leftValue == "") != (rightValue == "") {
		log.Fatal(fmt.Sprintf("%s and %s must be provided together", left, right))
	}
}
