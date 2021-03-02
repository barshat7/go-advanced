package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	ls "keyvalstore/localstore"
	l "keyvalstore/transactionlog"
	"fmt"
)

var transact l.TransactionLogger

func initializeTransactionLog() error {
	var err error 
	transact, err = l.New("transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create transaction")
	}
	events, errors := transact.ReadEvents()
	ok, e := true, l.Event{}
	for ok && err == nil {
		select {
		case err, ok = <- errors:
		case e, ok = <- events:
			switch e.EventType {
			case l.EventDelete:
				err = ls.Delete(e.Key)
			case l.EventPut:
				err = ls.Put(e.Key, e.Value)
			}
		}
	}
	transact.Run()
	return err
}

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
	transact.WritePut(key, string(value))
	w.WriteHeader(http.StatusCreated)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	err := ls.Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	transact.WriteDelete(key)
	w.WriteHeader(http.StatusAccepted)
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
	initializeTransactionLog()
	r := mux.NewRouter()
	r.HandleFunc("/", helloHandler)
	r.HandleFunc("/v1/{key}", putHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", getHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", deleteHandler).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":2181", r))
}