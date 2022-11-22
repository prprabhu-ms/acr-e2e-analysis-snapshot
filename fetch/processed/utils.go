package processed

import (
	"encoding/json"
	"log"
	"os"
)

func writeJSON(path string, data interface{}) {
	oData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(path, oData, 0755); err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %s\n", path)
}

func readJSON(path string, v interface{}) {
	t, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(t, v); err != nil {
		log.Fatal(err)
	}
}
