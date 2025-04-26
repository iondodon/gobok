package main

import (
	"fmt"
	"time"
)

func main() {
	// Create an address
	address := NewAddressBuilder().
		Street("123 Main St").
		City("New York").
		Country("USA").
		Build()

	// Create a contact with the address
	phone := "+1-555-1234"
	contact := NewContactBuilder().
		Email("john@example.com").
		Phone(&phone).
		Address(address).
		IsActive(true).
		Build()

	// Create settings map
	settings := map[string]string{
		"theme": "dark",
		"lang":  "en",
	}

	// Create metadata map
	metadata := map[string]interface{}{
		"preferences": []string{"email", "sms"},
		"last_login":  time.Now().Unix(),
	}

	// Create a user profile with all the above
	now := time.Now().Unix()
	profile := NewUserProfileBuilder().
		ID(1).
		Name("John Doe").
		Age(30).
		Contacts([]Contact{*contact}).
		Metadata(metadata).
		Settings(&settings).
		CreatedAt(now).
		UpdatedAt(&now).
		Build()

	// Print the results
	fmt.Println("Address:")
	fmt.Printf("%+v\n\n", address)

	fmt.Println("Contact:")
	fmt.Printf("%+v\n\n", contact)

	fmt.Println("User Profile:")
	fmt.Printf("%+v\n", profile)

	// Create some values for pointers
	boolVal := true
	intVal := 42
	stringVal := "pointer"

	// Create a nested struct
	nested := NewNestedStructBuilder().
		Field1("nested").
		Field2(123).
		Field3(&boolVal).
		Build()

	// Create channels
	intChan := make(chan int)
	sendChan := make(chan string)
	receiveChan := make(chan bool)

	// Create maps
	simpleMap := map[string]int{"one": 1, "two": 2}
	complexMap := map[string]map[int]string{
		"first":  {1: "one"},
		"second": {2: "two"},
	}
	interfaceMap := map[string]interface{}{
		"string": "value",
		"int":    42,
		"bool":   true,
	}
	structMap := map[string]NestedStruct{
		"first": *nested,
	}

	// Create the all types instance
	allTypes := NewAllTypesBuilder().
		// Basic types
		BoolValue(true).
		IntValue(42).
		Int8Value(8).
		Int16Value(16).
		Int32Value(32).
		Int64Value(64).
		UintValue(42).
		Uint8Value(8).
		Uint16Value(16).
		Uint32Value(32).
		Uint64Value(64).
		Float32Value(32.0).
		Float64Value(64.0).
		StringValue("string").
		ByteValue('A').
		RuneValue('世').

		// Pointer types
		BoolPtr(&boolVal).
		IntPtr(&intVal).
		StringPtr(&stringVal).
		StructPtr(nested).

		// Qualified types
		TimeValue(time.Now()).

		// Array types
		IntArray([]int{1, 2, 3}).
		StringArray([]string{"one", "two", "three"}).
		StructArray([]NestedStruct{*nested}).

		// Map types
		SimpleMap(simpleMap).
		ComplexMap(complexMap).
		InterfaceMap(interfaceMap).
		StructMap(structMap).

		// Channel types
		IntChan(intChan).
		SendChan(sendChan).
		ReceiveChan(receiveChan).

		// Nested struct
		NestedStruct(*nested).
		Build()

	// Print the results
	fmt.Printf("AllTypes: %+v\n", allTypes)
}
