package main

import (
	"net/http"
	"log"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

func main() {
	getJSON("http://localhost:3000/users")
}

func get(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Panicln("Could not Get")
	}
	fmt.Println(resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(string(body))
	resp.Body.Close()
}

type User struct {
	ID int
	Username string
}

func getJSON(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Panicln("Could not get")
	}
	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		log.Fatalln(err)
	}
	for _, user := range users {
		fmt.Println(user.ID, " ", user.Username)
	}
	defer resp.Body.Close()	
}

