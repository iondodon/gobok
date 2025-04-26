package main

import "time"

//gobok:builder
type AllTypes struct {
	// Basic types
	BoolValue    bool
	IntValue     int
	Int8Value    int8
	Int16Value   int16
	Int32Value   int32
	Int64Value   int64
	UintValue    uint
	Uint8Value   uint8
	Uint16Value  uint16
	Uint32Value  uint32
	Uint64Value  uint64
	Float32Value float32
	Float64Value float64
	StringValue  string
	ByteValue    byte
	RuneValue    rune

	// Pointer types
	BoolPtr   *bool
	IntPtr    *int
	StringPtr *string
	StructPtr *NestedStruct

	// Qualified types
	TimeValue time.Time

	// Array types
	IntArray    []int
	StringArray []string
	StructArray []NestedStruct

	// Map types
	SimpleMap    map[string]int
	ComplexMap   map[string]map[int]string
	InterfaceMap map[string]interface{}
	StructMap    map[string]NestedStruct

	// Channel types
	IntChan     chan int
	SendChan    chan<- string
	ReceiveChan <-chan bool

	// Nested struct
	NestedStruct NestedStruct
}

//gobok:builder
type NestedStruct struct {
	Field1 string
	Field2 int
	Field3 *bool
}
