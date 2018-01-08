package util

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func RpcCall(description string, f func() bool) {
	log.Println(description)

	t1 := time.Now().UnixNano()

	ok := f()
	if ok {
		cost := fmt.Sprintf("%.5f", (float64(time.Now().UnixNano()-t1) / float64(time.Second)))
		log.Print("[*] 成功, 用时" + cost + "秒")
	} else {
		log.Println("[*] 失败")
		log.Println("[*] 退出程序")
		os.Exit(0)
	}
}

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
