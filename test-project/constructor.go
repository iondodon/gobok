package main

//go:generate gobok

// Person represents a person with basic information
//
//gobok:builder
//gobok:constructor:name=CreatePerson
type Person struct {
	Name string
	Age  int
}

// Employee represents an employee with additional information
//
//gobok:builder
//gobok:constructor
type Employee struct {
	ID     int
	Title  string
	Salary float64
}
