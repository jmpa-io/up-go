package up

import (
	"encoding/json"
	"os"
)

// PLEASE NOTE:
// This file is intended to be used for reading & parsing json files from the
// testdata directory. Any code within this file is intended to be used within
// tests ONLY.

// readTestdata reads the testdata file at the given path and unmarshals it
// onto the given data (ideally a struct).
func readTestdata(path string, data interface{}) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &data)
}
