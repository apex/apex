package apex

import (
	"encoding/json"
	"io/ioutil"
)

// environment map, loaded from .env.json injected into the Lambda function's container.
var environ map[string]string

// Getenv retrieves the value of the environment variable named by the key.
func Getenv(name string) string {
	if environ == nil {
		b, err := ioutil.ReadFile(".env.json")
		if err != nil {
			return ""
		}

		if err := json.Unmarshal(b, &environ); err != nil {
			return ""
		}
	}

	return environ[name]
}
