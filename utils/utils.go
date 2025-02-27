package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

// PrettyPrint outputs a JSON-formatted representation of the given struct
func PrettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Printf("Error pretty printing: %v", err)
		return
	}
	fmt.Println(string(b))
}

// SafeJSONMarshal marshals an object to JSON and handles any errors
func SafeJSONMarshal(v interface{}) ([]byte, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("JSON marshal error: %w", err)
	}
	return bytes, nil
}
