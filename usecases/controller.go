package usecases

import (
	"SE-Publisher/model"
	"SE-Publisher/repository"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nsqio/go-nsq"
)

var db *gorm.DB
var err error

var mq *nsq.Producer

func init() {

	db, err = repository.CreateDbConnection()
	if err != nil {
		log.Fatal("error while creating connection :", err)
	}
	fmt.Println("successfully made connection", db)

	mq, err = repository.NewNsqConn()
	if err != nil {
		log.Fatal("error during connection: ", err)
	}
	fmt.Println("successfully made connection with messaging queue", mq)
	db.AutoMigrate(&model.Person{})
	db.AutoMigrate(&model.Book{})
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	var people []model.Person

	db.Find(&people)

	json.NewEncoder(w).Encode(&people)
	w.WriteHeader(http.StatusOK)

}

func GetPerson(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	// var books []model.Book
	var person model.Person
	db.First(&person, params["id"])
	// db.Model(&person).Related(&books)
	// if len(books) == 0 {
	// 	http.Error(w, "Person not present", http.StatusBadRequest)
	// 	return
	// }
	// person.Book = books
	json.NewEncoder(w).Encode(person)
}

func GetBook(w http.ResponseWriter, r *http.Request) {

	var books []model.Book
	db.Find(&books)
	json.NewEncoder(w).Encode(&books)

}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var person model.Person

	json.NewDecoder(r.Body).Decode(&person)

	createdPerson := db.Create(&person)
	err = createdPerson.Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		json.NewEncoder(w).Encode(&createdPerson)
	}
	topic := "test"
	pay, err := json.Marshal(model.Payload{Task: "created", Person: model.Person{Name: person.Name, Email: person.Email, Age: person.Age, DocumentId: person.DocumentId}})
	if err != nil {
		log.Fatal("error occured during parsing", err)
	}
	publishErr := mq.Publish(topic, pay)
	if publishErr != nil {
		log.Fatal(publishErr)
	}
	fmt.Println("payload is ", pay)
	fmt.Println("started listening events")

}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person model.Person

	db.First(&person, params["id"])
	db.Delete(&person)
	json.NewEncoder(w).Encode(&person)
	topic := "test"
	pay, err := json.Marshal(model.Payload{Task: "deleted", Person: model.Person{Name: person.Name, Email: person.Email, Age: person.Age, DocumentId: person.DocumentId}})
	if err != nil {
		log.Fatal("error occured during parsing", err)
	}
	publishErr := mq.Publish(topic, pay)
	if publishErr != nil {
		log.Fatal(publishErr)
	}
	fmt.Println("payload is ", pay)
	fmt.Println("deleted events sent successfully")
}

func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person model.Person

	db.First(&person, params["id"])

	json.NewDecoder(r.Body).Decode(&person)

	db.Save(&person)
	fmt.Println("updated successfully")
	topic := "test"
	pay, err := json.Marshal(model.Payload{Task: "updated", Person: model.Person{Name: person.Name, Email: person.Email, Age: person.Age, DocumentId: person.DocumentId}})
	if err != nil {
		log.Fatal("error occured during parsing", err)
	}
	publishErr := mq.Publish(topic, pay)
	if publishErr != nil {
		log.Fatal(publishErr)
	}
	fmt.Println("payload is ", pay)
	fmt.Println("Updated events sent successfully")

}

func DeleteAll(w http.ResponseWriter, r *http.Request) {

	var person []model.Person

	db.Delete(&person)
}
