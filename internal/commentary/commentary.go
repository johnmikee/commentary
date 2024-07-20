package commentary

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ProcessDirectory processes all .go files in the given directory.
func ProcessDirectory(dir string, write bool) error {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".go") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking the path: %w", err)
	}

	for _, file := range files {
		fmt.Printf("Processing file: %s\n", file)
		processFile(file, write)
	}
	return nil
}

func processFile(filename string, write bool) {
	source, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", filename, err)
		return
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, source, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing file:", filename, err)
		return
	}

	var builder strings.Builder
	lastPos := 0
	edited := false

	for _, commentGroup := range node.Comments {
		for _, comment := range commentGroup.List {
			if shouldSkipComment(comment.Text) {
				fmt.Printf("Skipping special comment: %s\n", comment.Text)
				continue
			}
			newText := processComment(comment, node)
			if newText != comment.Text {
				fmt.Printf("Modifying comment: %s -> %s\n", comment.Text, newText)
				edited = true
			}
			start := fset.Position(comment.Pos()).Offset
			end := fset.Position(comment.End()).Offset
			builder.Write(source[lastPos:start])
			builder.WriteString(newText)
			lastPos = end
		}
	}
	builder.Write(source[lastPos:])

	if edited {
		if write {
			err := os.WriteFile(filename, []byte(builder.String()), 0o644)
			if err != nil {
				fmt.Println("Error writing file:", filename, err)
			} else {
				fmt.Printf("Changes written to file: %s\n", filename)
			}
		} else {
			fmt.Printf("Changes detected but not written to file: %s\n", filename)
		}
	} else {
		fmt.Printf("No changes needed for file: %s\n", filename)
	}
}

func shouldSkipComment(text string) bool {
	return strings.HasPrefix(text, "/*") || strings.HasPrefix(text, "//nolint") || strings.HasPrefix(text, "//go:") || strings.HasPrefix(text, "//+build")
}

func processComment(comment *ast.Comment, node *ast.File) string {
	if isExportedComment(comment, node) {
		return capitalizeFirst(comment.Text)
	}
	return lowercaseFirst(comment.Text)
}

func isExportedComment(comment *ast.Comment, node *ast.File) bool {
	commentPos := comment.Slash
	for _, decl := range node.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.TypeSpec:
					if decl.Doc != nil && commentPos == decl.Doc.Pos() {
						return spec.Name.IsExported()
					}
					if spec.Doc != nil && commentPos == spec.Doc.Pos() {
						return spec.Name.IsExported()
					}
					if structType, ok := spec.Type.(*ast.StructType); ok {
						for _, field := range structType.Fields.List {
							if field.Doc != nil && commentPos == field.Doc.Pos() {
								return field.Names[0].IsExported()
							}
						}
					}
				}
			}
		case *ast.FuncDecl:
			if decl.Doc != nil && commentPos == decl.Doc.Pos() {
				return decl.Name.IsExported()
			}
			if decl.Body != nil {
				for _, stmt := range decl.Body.List {
					if hasInternalComment(stmt, commentPos) {
						return false
					}
				}
			}
		}
	}
	return false
}

func hasInternalComment(n ast.Node, pos token.Pos) bool {
	var found bool
	ast.Inspect(n, func(n ast.Node) bool {
		if found {
			return false
		}
		if cg, ok := n.(*ast.CommentGroup); ok {
			for _, c := range cg.List {
				if c.Pos() == pos {
					found = true
					return false
				}
			}
		}
		return true
	})
	return found
}

func capitalizeFirst(s string) string {
	if len(s) <= 3 {
		return s
	}
	return "// " + strings.ToUpper(string(s[3])) + s[4:]
}

func lowercaseFirst(s string) string {
	if len(s) <= 3 {
		return s
	}
	return "// " + strings.ToLower(string(s[3])) + s[4:]
}
