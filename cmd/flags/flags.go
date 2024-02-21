package flags

import (
	"fmt"

	"github.com/spf13/viper"
)

func GetString(key string) (string, error) {
	value := viper.GetString(key)
	if value == "" {
		return "", fmt.Errorf("flag '%s' is required", key)
	}
	return value, nil
}
