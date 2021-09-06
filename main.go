package main

import (
	"net/http"

	"SE-Publisher/usecases"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/people", usecases.GetPeople).Methods("GET")
	router.HandleFunc("/books", usecases.GetBook).Methods("GET")
	router.HandleFunc("/person/{id}", usecases.GetPerson).Methods("GET")
	router.HandleFunc("/create/person", usecases.CreatePerson).Methods("POST")
	router.HandleFunc("/delete/{id}", usecases.DeletePerson).Methods("DELETE")
	router.HandleFunc("/update/{id}", usecases.UpdatePerson).Methods("PUT")
	router.HandleFunc("/deleteAll", usecases.DeleteAll).Methods("DELETE")
	http.ListenAndServe(":8080", router)

}
