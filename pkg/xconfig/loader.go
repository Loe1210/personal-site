package xconfig

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

func LoadYAML(path string, target interface{}) error {
	if path == "" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return yaml.Unmarshal(data, target)
}

func OverrideString(dst *string, envName string) {
	if value := os.Getenv(envName); value != "" {
		*dst = value
	}
}

func OverrideInt(dst *int, envName string) {
	if value := os.Getenv(envName); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			*dst = parsed
		}
	}
}

func OverrideInt64(dst *int64, envName string) {
	if value := os.Getenv(envName); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			*dst = parsed
		}
	}
}
