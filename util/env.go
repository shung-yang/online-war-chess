package util

import (
	"fmt"
	"os"
)

func ReadEnvVariable(key string) (string, bool) {
	value, not_empty := os.LookupEnv(key)
	if !not_empty {
		fmt.Printf("%s not set\n", key)
		return "", false
	}
	return value, true
}