package sharedcrud

import (
	"strings"

	shared "github.com/max38/golang-clean-code-architecture/src/shared"
)

func ConvertNameToCRUDSlug(name string) string {
	// camelCase to snake_case
	name = shared.CamelToSnake(name)

	return strings.Replace(name, "_", "-", -1)
}
