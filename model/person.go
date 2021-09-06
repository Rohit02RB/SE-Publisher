package model

import "github.com/jinzhu/gorm"

type Person struct {
	gorm.Model

	Name       string
	Email      string
	Age        int
	Password   int
	Book       []Book
	DocumentId int
}

type Book struct {
	gorm.Model

	Title    string
	Author   string
	PersonId int
}

type Payload struct {
	Task   string
	Person Person
}
