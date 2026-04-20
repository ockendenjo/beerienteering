package env

import "os"

func OptStr(key, defaultValue string) string {
	s := os.Getenv(key)
	if s == "" {
		return defaultValue
	}
	return s
}
