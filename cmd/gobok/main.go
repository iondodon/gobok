package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const toolVersion = "v1.0.0"

type BuilderData struct {
	StructName          string
	Fields              []FieldData
	GenerateBuilder     bool
	GenerateConstructor bool
	ConstructorName     string
}

type FieldData struct {
	Name string
	Type string
}

type FolderData struct {
	PackageName string
	Builders    []BuilderData
	Imports     map[string]string // Track required imports with their full paths
}

var folders = make(map[string]*FolderData)

type ImportData struct {
	Alias string
	Path  string
}

type TemplateData struct {
	PackageName string
	Builders    []BuilderData
	ToolVersion string
	Imports     []ImportData
}

func main() {
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				base := filepath.Base(path)
				if base == "vendor" || base == ".git" || strings.HasPrefix(base, ".") {
					return filepath.SkipDir
				}
				return nil
			}

			if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") || filepath.Base(path) == "gobok.go" {
				return nil
			}

			processFile(path)
			return nil
		})

		if err != nil {
			fmt.Printf("Error walking %s: %v\n", root, err)
		}
	}

	for folder, data := range folders {
		writeBuilders(folder, data)
	}
}

func processFile(path string) {
	fmt.Printf("[gobok] Scanning file: %s\n", path)

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Failed to parse %s: %v\n", path, err)
		return
	}

	folder := filepath.Dir(path)
	if folders[folder] == nil {
		folders[folder] = &FolderData{
			PackageName: node.Name.Name,
			Imports:     make(map[string]string),
		}
	}

	// Create a map of original imports for reference
	originalImports := make(map[string]string)
	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, "\"")
		if imp.Name != nil {
			originalImports[imp.Name.Name] = importPath
		} else {
			parts := strings.Split(importPath, "/")
			originalImports[parts[len(parts)-1]] = importPath
		}
	}

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		if genDecl.Doc == nil {
			continue
		}

		builder := BuilderData{}

		for _, comment := range genDecl.Doc.List {
			text := strings.TrimSpace(comment.Text)

			switch {
			case text == "//gobok:builder":
				builder.GenerateBuilder = true
			case text == "//gobok:constructor":
				builder.GenerateConstructor = true
			case strings.HasPrefix(text, "//gobok:constructor:name="):
				builder.GenerateConstructor = true
				builder.ConstructorName = strings.TrimPrefix(text, "//gobok:constructor:name=")
			}
		}

		if !builder.GenerateBuilder && !builder.GenerateConstructor {
			continue
		}

		typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec)
		if !ok {
			continue
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		builder.StructName = typeSpec.Name.Name

		for _, field := range structType.Fields.List {
			for _, name := range field.Names {
				fieldType := exprToString(field.Type)
				builder.Fields = append(builder.Fields, FieldData{
					Name: name.Name,
					Type: fieldType,
				})

				// Track imports from field types
				if strings.Contains(fieldType, ".") {
					parts := strings.Split(fieldType, ".")
					if len(parts) > 1 {
						pkg := parts[0]
						// Handle pointer types
						if strings.HasPrefix(pkg, "*") {
							pkg = pkg[1:]
						}
						// Handle array types
						if strings.HasPrefix(pkg, "[]") {
							pkg = pkg[2:]
						}
						// Handle map types
						if strings.HasPrefix(pkg, "map[") {
							pkg = strings.TrimPrefix(pkg, "map[")
							pkg = strings.Split(pkg, "]")[0]
						}
						// Add the import if it's not a built-in type
						if !isBuiltInType(pkg) {
							// Try to find the import in the original imports
							if importPath, exists := originalImports[pkg]; exists {
								folders[folder].Imports[pkg] = importPath
							}
						}
					}
				}
			}
		}

		folders[folder].Builders = append(folders[folder].Builders, builder)
	}
}

// Helper function to check if a type is a built-in Go type
func isBuiltInType(typeName string) bool {
	builtInTypes := map[string]bool{
		"bool":       true,
		"string":     true,
		"int":        true,
		"int8":       true,
		"int16":      true,
		"int32":      true,
		"int64":      true,
		"uint":       true,
		"uint8":      true,
		"uint16":     true,
		"uint32":     true,
		"uint64":     true,
		"uintptr":    true,
		"byte":       true,
		"rune":       true,
		"float32":    true,
		"float64":    true,
		"complex64":  true,
		"complex128": true,
		"error":      true,
		"interface":  true,
	}
	return builtInTypes[typeName]
}

func writeBuilders(folder string, data *FolderData) {
	tmpl, err := template.New("builder").Parse(builderTemplate)
	if err != nil {
		fmt.Printf("Failed to parse template: %v\n", err)
		return
	}

	// Convert imports map to slice of ImportData
	imports := make([]ImportData, 0, len(data.Imports))
	for alias, path := range data.Imports {
		// Extract the last part of the path
		parts := strings.Split(path, "/")
		lastPart := parts[len(parts)-1]

		// If the alias matches the last part of the path, we don't need an alias
		if alias == lastPart {
			imports = append(imports, ImportData{
				Path: path,
			})
		} else {
			imports = append(imports, ImportData{
				Alias: alias,
				Path:  path,
			})
		}
	}

	outData := TemplateData{
		PackageName: data.PackageName,
		Builders:    data.Builders,
		ToolVersion: toolVersion,
		Imports:     imports,
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, outData)
	if err != nil {
		fmt.Printf("Failed to execute template: %v\n", err)
		return
	}

	source, err := format.Source([]byte(buf.String()))
	if err != nil {
		fmt.Printf("Failed to format generated code: %v\n", err)
		source = []byte(buf.String())
	}

	outPath := filepath.Join(folder, "gobok.go")
	fmt.Printf("[gobok] Generating file: %s\n", outPath)
	err = os.WriteFile(outPath, source, 0644)
	if err != nil {
		fmt.Printf("Failed to write file %s: %v\n", outPath, err)
	}
}

func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + exprToString(t.X)
	case *ast.SelectorExpr:
		return exprToString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + exprToString(t.Elt)
	case *ast.MapType:
		return "map[" + exprToString(t.Key) + "]" + exprToString(t.Value)
	case *ast.ChanType:
		dir := ""
		if t.Dir == ast.SEND {
			dir = "chan<- "
		} else if t.Dir == ast.RECV {
			dir = "<-chan "
		} else {
			dir = "chan "
		}
		return dir + exprToString(t.Value)
	case *ast.FuncType:
		var params []string
		for _, f := range t.Params.List {
			for range f.Names {
				params = append(params, exprToString(f.Type))
			}
		}

		var results []string
		if t.Results != nil {
			for _, f := range t.Results.List {
				for _ = range f.Names {
					results = append(results, exprToString(f.Type))
				}
				if len(f.Names) == 0 {
					results = append(results, exprToString(f.Type))
				}
			}
		}

		paramList := strings.Join(params, ", ")
		resultList := strings.Join(results, ", ")

		if len(results) == 1 {
			return fmt.Sprintf("func(%s) %s", paramList, resultList)
		} else if len(results) > 1 {
			return fmt.Sprintf("func(%s) (%s)", paramList, resultList)
		}
		return fmt.Sprintf("func(%s)", paramList)

	case *ast.StructType:
		var fields []string
		for _, f := range t.Fields.List {
			for _, name := range f.Names {
				fields = append(fields, fmt.Sprintf("%s %s", name.Name, exprToString(f.Type)))
			}
			if len(f.Names) == 0 {
				fields = append(fields, exprToString(f.Type))
			}
		}
		return fmt.Sprintf("struct { %s }", strings.Join(fields, "; "))

	case *ast.InterfaceType:
		var methods []string
		for _, f := range t.Methods.List {
			for _, name := range f.Names {
				methods = append(methods, fmt.Sprintf("%s %s", name.Name, exprToString(f.Type)))
			}
			if len(f.Names) == 0 {
				methods = append(methods, exprToString(f.Type))
			}
		}
		return fmt.Sprintf("interface { %s }", strings.Join(methods, "; "))

	default:
		return "interface{}"
	}
}
