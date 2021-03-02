package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	ls "keyvalstore/localstore"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Hello Handler Executing...")
	w.Write([] byte("Hello net/http\n"))
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}
	err = ls.Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value, err := ls.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([] byte(value))
}


func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", helloHandler)
	r.HandleFunc("/v1/{key}", putHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", getHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":2181", r))
}