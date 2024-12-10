package gassfile

import "strings"

func stringInSlice(str string, sl []string) bool {
	for _, s := range sl {
		if strings.EqualFold(s, str) {
			return true
		}
	}

	return false
}
