package main

//gobok:builder
type Person struct {
	Name     string
	Age      int
	IsActive bool
	Tags     []string
	Metadata map[string]interface{}
	Parent   *Person
}
