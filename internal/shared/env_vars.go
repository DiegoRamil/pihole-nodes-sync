package shared

import "os"

func RetrieveEnvVar(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("make sure that you update the environment variables")
	}
	return v
}
