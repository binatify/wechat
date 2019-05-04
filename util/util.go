package util

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

func UnixTimestamp() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

func JsonDecode(content string) interface{} {
	content = strings.Replace(content, "\n", "", -1)

	var ret interface{}

	err := json.Unmarshal([]byte(content), &ret)
	if err != nil {
		return false
	}

	return formatData(ret)
}

func formatData(input interface{}) interface{} {
	if m, ok := input.([]interface{}); ok {
		for k, v := range m {
			switch v.(type) {
			case float64:
				m[k] = int(v.(float64))
			case []interface{}:
				m[k] = formatData(m[k])
			case map[string]interface{}:
				m[k] = formatData(m[k])
			}
		}
	} else if m, ok := input.(map[string]interface{}); ok {
		for k, v := range m {
			switch v.(type) {
			case float64:
				m[k] = int(v.(float64))
			case []interface{}:
				m[k] = formatData(m[k])
			case map[string]interface{}:
				m[k] = formatData(m[k])
			}
		}
	} else {
		return false
	}

	return input
}
