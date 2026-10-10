package shared

import "encoding/json"

func PrettyPrintJson(data interface{}) string {
	json, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return ""
	}
	return string(json)
}
