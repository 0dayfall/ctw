package utils

import (
	"encoding/json"
	"log"
)

func PrettyPrint(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		log.Println(string(b))
	}
}

func FindString(data any, target string) bool {
	switch v := data.(type) {
	case string:
		return v == target
	case []any:
		for _, elem := range v {
			if FindString(elem, target) {
				return true
			}
		}
	case map[string]any:
		for _, elem := range v {
			if FindString(elem, target) {
				return true
			}
		}
	}
	return false
}

func FindStringWithValue(data any, target string) bool {
	switch v := data.(type) {
	case string:
		return v == target
	case []any:
		for _, elem := range v {
			if FindString(elem, target) {
				return true
			}
		}
	case map[string]any:
		for _, elem := range v {
			if FindString(elem, target) {
				return true
			}
		}
	}
	return false
}
