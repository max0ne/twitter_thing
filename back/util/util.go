package util

import (
	"encoding/json"
	"log"
	"os"
)

// Contains - -
func Contains(vals []string, aVal string) bool {
	for _, v := range vals {
		if v == aVal {
			return true
		}
	}
	return false
}

// Remove - -
func Remove(vals []string, aVal string) []string {
	newVals := []string{}
	for _, v := range vals {
		if v == aVal {
			continue
		}
		newVals = append(newVals, v)
	}
	return newVals
}

// JSONMarshel - -
func JSONMarshel(val interface{}) string {
	str, _ := json.Marshal(val)
	return string(str)
}

// GetEnvMust get env, crashes if env key not set
func GetEnvMust(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatal("env key ", key, " missing")
	}
	return val
}
