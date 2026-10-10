package shared

import (
	"bytes"
	"strconv"
	"unicode"
)

func ConvertStringToInteger(stringValue string) int {
	integerValue, err := strconv.Atoi(stringValue)
	if err != nil {
		// show log
		return 0
	}
	return integerValue
}

func CamelToSnake(camel string) string {
	var snake bytes.Buffer
	for i, r := range camel {
		if unicode.IsUpper(r) {
			if i > 0 {
				snake.WriteRune('_')
			}
			snake.WriteRune(unicode.ToLower(r))
		} else {
			snake.WriteRune(r)
		}
	}
	return snake.String()
}
