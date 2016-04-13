package main

import (
	"os"
	"strconv"
)

var PORT_ENV_PREFIX = "AUTO_PORT"

func GetEnv(name string) string {
	return os.Getenv(name)
}

func SetEnv(name string, value string) {
	os.Setenv(name, value)
}

func FormatPortEnv(ports []int) map[string]string {
	portEnv := make(map[string]string)
	for i, port := range ports {
		portEnv[PORT_ENV_PREFIX+strconv.Itoa(i)] = strconv.Itoa(port)
	}
	return portEnv
}

func ExportEnvs(ports map[string]string) {
	for key, value := range ports {
		SetEnv(key, value)
	}
}
