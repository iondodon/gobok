# gobok

gobok is a Go code generator that automatically generates builder patterns and constructors for your structs. It helps reduce boilerplate code while providing a fluent interface for object construction.

## Features

- **Builder Pattern Generation**: Creates fluent builder interfaces for your structs
- **Constructor Generation**: Creates all-args constructors for your structs
- **Concise Method Names**: Uses field names directly as method names (e.g., `Name()` instead of `SetName()`)
- **Support for All Go Types**: Works with basic types, pointers, arrays, maps, channels, and custom types
- **Custom Constructor Names**: Allows specifying custom names for constructors

## Installation

```bash
go install github.com/iondodon/gobok/cmd/gobok@latest
```

## Usage

### 1. Add Directives to Your Structs

Add `//gobok:` directives above your struct definitions:

```go
// Generate only a builder
//gobok:builder
type Person struct {
    Name string
    Age  int
}

// Generate both a builder and a constructor with default name
//gobok:builder
//gobok:constructor
type Employee struct {
    ID     int
    Title  string
    Salary float64
}

// Generate both a builder and a constructor with custom name
//gobok:builder
//gobok:constructor:name=CreateUser
type User struct {
    Name  string
    Email string
}
```

### 2. Generate Code

Run gobok on your project:

```bash
gobok .     # Process project root and all subdirectories recursively
gobok ./pkg # Process specific package
```

### 3. Use Generated Code

```go
// Using a builder
person := NewPersonBuilder().
    Name("John").
    Age(30).
    Build()

// Using a default constructor
employee := NewEmployee(1, "Developer", 75000.0)

// Using a custom-named constructor
user := CreateUser("Alice", "alice@example.com")
```

## Generated Code Example

For a struct like:

```go
//gobok:builder
//gobok:constructor
type Person struct {
    Name string
    Age  int
}
```

gobok generates:

```go
// Builder
type PersonBuilder struct {
    instance *Person
}

func NewPersonBuilder() *PersonBuilder {
    return &PersonBuilder{
        instance: &Person{},
    }
}

func (b *PersonBuilder) Name(v string) *PersonBuilder {
    b.instance.Name = v
    return b
}

func (b *PersonBuilder) Age(v int) *PersonBuilder {
    b.instance.Age = v
    return b
}

func (b *PersonBuilder) Build() *Person {
    return b.instance
}

// Constructor
func NewPerson(Name string, Age int) Person {
    return Person{
        Name: Name,
        Age:  Age,
    }
}
```

## Directives

- `//gobok:builder`: Generates a builder for the struct
- `//gobok:constructor`: Generates a constructor with default name (New[StructName])
- `//gobok:constructor:name=CustomName`: Generates a constructor with a custom name

## License

MIT License
