package util

import (
	"encoding/json"
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
