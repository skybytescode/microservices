package config

import (
	"log"
	"os"
	"strconv"
)

func GetEnv() string {
	return getEnvironmentValue("ENV")
}

func GetDataSourceURL() string {
	return getEnvironmentValue("DATA_SOURCE_URL")
}

func GetApplicationPort() int {
	portStr := getEnvironmentValue("APPLICATION_PORT")
	port, err := strconv.Atoi(portStr)

	if err != nil {
		log.Fatalf("port: %s is invalid", portStr)
	}

	return port
}

func GetPaymentServiceUrl() string {
	return getEnvironmentValue("PAYMENT_SERVICE_URL")
}

func GetShippingServiceUrl() string {
	return getEnvironmentValue("SHIPPING_SERVICE_URL")
}

// GetMetricsPort returns the port that serves Prometheus metrics, or 0 when
// METRICS_PORT is not set, which turns the metrics endpoint off.
func GetMetricsPort() int {
	portStr := os.Getenv("METRICS_PORT")
	if portStr == "" {
		return 0
	}
	port, err := strconv.Atoi(portStr)

	if err != nil {
		log.Fatalf("metrics port: %s is invalid", portStr)
	}

	return port
}

func getEnvironmentValue(key string) string {
	if os.Getenv(key) == "" {
		log.Fatalf("%s environment variable is missing.", key)
	}

	return os.Getenv(key)
}
