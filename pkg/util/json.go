package util

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ToJson converts an object to JSON string.
// Returns empty string if marshaling fails.
// Deprecated: Use ToJsonSafe for proper error handling.
func ToJson(obj interface{}) string {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

// ToJsonSafe converts an object to JSON string with error handling.
// Returns an error if marshaling fails, allowing callers to handle failures appropriately.
func ToJsonSafe(obj interface{}) (string, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("json marshal failed: %w", err)
	}
	return string(jsonData), nil
}

func JsonToMap(json string) (map[string]interface{}, error) {
	var resultMap map[string]interface{}
	err := ReadJsonByByte([]byte(json), &resultMap)
	return resultMap, err
}
func ReadJsonByByte(body []byte, obj interface{}) error {
	mdz := json.NewDecoder(bytes.NewBuffer(body))

	mdz.UseNumber()
	err := mdz.Decode(obj)

	if err != nil {
		return err
	}
	return nil
}
