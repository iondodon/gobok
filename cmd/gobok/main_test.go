package main

import (
	"go/parser"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProcessFile(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.go")

	content := `package test

//gobok:builder
type TestStruct struct {
	Name string
	Age  int
	Tags []string
}`

	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Process the test file
	processFile(testFile)

	// Verify the folder data was created
	folder := tempDir
	if folders[folder] == nil {
		t.Fatal("Folder data was not created")
	}

	// Verify the builder data
	if len(folders[folder].Builders) != 1 {
		t.Fatalf("Expected 1 builder, got %d", len(folders[folder].Builders))
	}

	builder := folders[folder].Builders[0]
	if builder.StructName != "TestStruct" {
		t.Errorf("Expected struct name 'TestStruct', got '%s'", builder.StructName)
	}

	if !builder.GenerateBuilder {
		t.Error("Expected GenerateBuilder to be true")
	}

	// Verify fields
	expectedFields := []FieldData{
		{Name: "Name", Type: "string"},
		{Name: "Age", Type: "int"},
		{Name: "Tags", Type: "[]string"},
	}

	if len(builder.Fields) != len(expectedFields) {
		t.Fatalf("Expected %d fields, got %d", len(expectedFields), len(builder.Fields))
	}

	for i, field := range builder.Fields {
		if field != expectedFields[i] {
			t.Errorf("Field %d mismatch: expected %v, got %v", i, expectedFields[i], field)
		}
	}
}

func TestExprToString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic type",
			input:    "string",
			expected: "string",
		},
		{
			name:     "pointer type",
			input:    "*int",
			expected: "*int",
		},
		{
			name:     "array type",
			input:    "[]string",
			expected: "[]string",
		},
		{
			name:     "map type",
			input:    "map[string]int",
			expected: "map[string]int",
		},
		{
			name:     "channel type",
			input:    "chan int",
			expected: "chan int",
		},
		{
			name:     "send channel",
			input:    "chan<- int",
			expected: "chan<- int",
		},
		{
			name:     "receive channel",
			input:    "<-chan int",
			expected: "<-chan int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tt.input)
			if err != nil {
				t.Fatalf("Failed to parse expression: %v", err)
			}

			result := exprToString(expr)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestWriteBuilders(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.go")

	content := `package test

//gobok:builder
type TestStruct struct {
	Name string
	Age  int
}`

	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Process the test file
	processFile(testFile)

	// Write the builders
	writeBuilders(tempDir, folders[tempDir])

	// Verify the generated file exists
	generatedFile := filepath.Join(tempDir, "gobok.go")
	if _, err := os.Stat(generatedFile); os.IsNotExist(err) {
		t.Fatal("Generated file was not created")
	}

	// Read and verify the generated content
	generatedContent, err := os.ReadFile(generatedFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	// Basic content verification
	contentStr := string(generatedContent)
	if !strings.Contains(contentStr, "type TestStructBuilder struct") {
		t.Error("Generated file does not contain builder struct")
	}
	if !strings.Contains(contentStr, "func NewTestStructBuilder() *TestStructBuilder") {
		t.Error("Generated file does not contain builder constructor")
	}
	if !strings.Contains(contentStr, "func (b *TestStructBuilder) Build() *TestStruct") {
		t.Error("Generated file does not contain Build method")
	}
}
