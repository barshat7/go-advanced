package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

/**
1. Read all the keys  of db.json
2. Create routes for each key
**/
func main() { 
	jsonFile, err := os.Open("db.json")
	defer jsonFile.Close()
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface {}
	json.Unmarshal([]byte(byteValue), &result)
	fmt.Println(result["users"])
}	