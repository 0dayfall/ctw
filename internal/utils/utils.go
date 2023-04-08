package utils

import (
	"encoding/json"
	"log"
	"strings"
)

func PrettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		log.Println(string(b))
	}
}

func FindString(data interface{}, target string) bool {
	switch v := data.(type) {
	case string:
		return strings.Compare(v, target) == 0
	case []interface{}:
		for _, elem := range v {
			if FindString(elem, target) {
				return true
			}
		}
	case map[string]interface{}:
		for _, elem := range v {
			if FindString(elem, target) {
				return true
			}
		}
	}
	return false
}

func FindStringWithValue(data interface{}, target string) bool {
	switch v := data.(type) {
	case string:
		return strings.Compare(v, target) == 0
	case []interface{}:
		for _, elem := range v {
			if FindString(elem, target) {
				return true
			}
		}
	case map[string]interface{}:
		for _, elem := range v {
			if FindString(elem, target) {
				return true
			}
		}
	}
	return false
}
