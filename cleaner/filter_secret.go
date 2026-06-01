package cleaner

import "encoding/base64"

// DecodeSecret 将 Secret 的 data 字段从 base64 解码，移到 stringData
func DecodeSecret(resource map[string]interface{}) map[string]interface{} {
	if kind, ok := resource["kind"].(string); !ok || kind != "Secret" {
		return resource
	}

	data, ok := resource["data"].(map[string]interface{})
	if !ok {
		return resource
	}

	stringData := make(map[string]interface{})
	for k, v := range data {
		if encoded, ok := v.(string); ok {
			if decoded, err := base64.StdEncoding.DecodeString(encoded); err == nil {
				stringData[k] = string(decoded)
			} else {
				stringData[k] = encoded
			}
		} else {
			stringData[k] = v
		}
	}

	delete(resource, "data")
	if len(stringData) > 0 {
		resource["stringData"] = stringData
	}

	return resource
}
