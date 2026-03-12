package util

import (
	"bytes"
	"log"
	"encoding/json"
)

//将对象转换为JSON
func ToJson(obj interface{}) string {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		log.Printf("ToJson marshal error: %v", err)
		return ""
	}
	return string(jsonData)
}

// ToJsonE 将对象转换为JSON，返回错误而非静默忽略
func ToJsonE(obj interface{}) (string, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return "", err
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
