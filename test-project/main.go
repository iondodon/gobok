package main

import (
	"fmt"
	"time"
)

func main() {
	// Create an address
	address := NewAddressBuilder().
		SetStreet("123 Main St").
		SetCity("New York").
		SetCountry("USA").
		Build()

	// Create a contact with the address
	phone := "+1-555-1234"
	contact := NewContactBuilder().
		SetEmail("john@example.com").
		SetPhone(&phone).
		SetAddress(address).
		SetIsActive(true).
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
		SetID(1).
		SetName("John Doe").
		SetAge(30).
		SetContacts([]Contact{*contact}).
		SetMetadata(metadata).
		SetSettings(&settings).
		SetCreatedAt(now).
		SetUpdatedAt(&now).
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
		SetField1("nested").
		SetField2(123).
		SetField3(&boolVal).
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
		SetBoolValue(true).
		SetIntValue(42).
		SetInt8Value(8).
		SetInt16Value(16).
		SetInt32Value(32).
		SetInt64Value(64).
		SetUintValue(42).
		SetUint8Value(8).
		SetUint16Value(16).
		SetUint32Value(32).
		SetUint64Value(64).
		SetFloat32Value(32.0).
		SetFloat64Value(64.0).
		SetStringValue("string").
		SetByteValue('A').
		SetRuneValue('世').

		// Pointer types
		SetBoolPtr(&boolVal).
		SetIntPtr(&intVal).
		SetStringPtr(&stringVal).
		SetStructPtr(nested).

		// Qualified types
		SetTimeValue(time.Now()).

		// Array types
		SetIntArray([]int{1, 2, 3}).
		SetStringArray([]string{"one", "two", "three"}).
		SetStructArray([]NestedStruct{*nested}).

		// Map types
		SetSimpleMap(simpleMap).
		SetComplexMap(complexMap).
		SetInterfaceMap(interfaceMap).
		SetStructMap(structMap).

		// Channel types
		SetIntChan(intChan).
		SetSendChan(sendChan).
		SetReceiveChan(receiveChan).

		// Nested struct
		SetNestedStruct(*nested).
		Build()

	// Print the results
	fmt.Printf("AllTypes: %+v\n", allTypes)
}
