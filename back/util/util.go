package util

func Contains(vals []string, aVal string) bool {
	for _, v := range vals {
		if v == aVal {
			return true
		}
	}
	return false
}

func Remove(vals []string, aVal string) []string {
	newVals := []string{}
	for _, v := range vals {
		if v == aVal {
			continue
		}
		newVals = append(newVals, aVal)
	}
	return newVals
}
