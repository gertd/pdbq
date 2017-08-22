package helper

import (
	"encoding/json"
	"log"
)

// PrettyPrintJSON : return a format JSON string representation
func PrettyPrintJSON(p interface{}) string {

	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		log.Println("error:", err)
		return ""
	}
	return string(b)
}
