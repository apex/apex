package apex

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func init() {
	env := make(map[string]string)

	b, err := ioutil.ReadFile(".env.json")
	if err != nil {
		return
	}

	if err := json.Unmarshal(b, &env); err != nil {
		return
	}

	for k, v := range env {
		os.Setenv(k, v)
	}
}
