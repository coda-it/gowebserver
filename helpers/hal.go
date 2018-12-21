package helpers

import (
	"encoding/json"
)

// ServeHal - serves HAL compliant responses
func ServeHal(data interface{}, embedded interface{}, links interface {}) map[string]interface{} {
	var res map[string]interface{}

	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &res)

	res["_links"] = links
	res["_embedded"] = embedded

	return res
}
