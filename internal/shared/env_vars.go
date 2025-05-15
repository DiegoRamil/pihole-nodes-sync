package shared

import (
	"fmt"
	"os"
)

func RetrieveEnvVar(key string) string {
	if os.Getenv("PROFILE") == "debug" || os.Getenv("PROFILE") == "DEBUG" {
		fmt.Printf("Retrieving environment variable: %s\n", key)
	}
	v := os.Getenv(key)
	if v == "" {
		panic("make sure that you update the environment variables")
	}
	return v
}
