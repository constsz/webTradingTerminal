package system

import (
	"encoding/json"
	"os"
)

type Credentials struct {
	ApiKey    string
	SecretKey string
}

// ReadJsonCredentials: Read `credentials.json` file for ApiKey and ApiSecret.
func ReadJsonCredentials() (Credentials, error) {
	// Open JSON file
	file, err := os.Open("credentials.json")
	if err != nil {
		return Credentials{}, err
	}
	defer file.Close()

	// Unmarshall json into struct
	var c Credentials
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		return Credentials{}, err
	}

	return c, nil

}
