package openapi

import "strings"

func camelfy(str string) string {
	transformed := ""

	for i := range str {
		letter := string(str[i])

		if letter != "_" {
			if i == 0 {
				transformed += strings.ToLower(letter)
				continue
			}

			transformed += letter
			continue
		}

		next := string(str[i+1])
		transformed += strings.ToUpper(next)
		i++
	}

	return transformed
}
