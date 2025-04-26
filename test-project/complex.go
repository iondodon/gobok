package main

//gobok:builder
type Address struct {
	Street  string
	City    string
	Country string
}

//gobok:builder
type Contact struct {
	Email    string
	Phone    *string
	Address  *Address
	IsActive bool
}

//gobok:builder
type UserProfile struct {
	ID        int
	Name      string
	Age       int
	Contacts  []Contact
	Metadata  map[string]interface{}
	Settings  *map[string]string
	CreatedAt int64
	UpdatedAt *int64
}
