package gosf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var Config map[string]interface{}

func init() {
	Config = make(map[string]interface{})
}

func LoadConfig(name string) {
	jsonFile, err := os.Open(name + ".json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var data map[string]interface{}
	json.Unmarshal([]byte(byteValue), &data)

	Config[name] = data
}
