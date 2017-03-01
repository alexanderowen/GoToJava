// Program reads a small subset of Go programs, parses it, and
// then prints the Java source equivalent
//
// Usage:
// 'go build main.go'
// './main filename'

package main

import (
	jv "JavaVisitor"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
)

var debug = false

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s filename\n", os.Args[0])
		os.Exit(-1)
	}

	fn := os.Args[1]
	file, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file\n")
		os.Exit(-1)
	}

	// Create the AST by parsing src.
	fset := token.NewFileSet() // token positions are relative to fset
	f, err := parser.ParseFile(fset, fn, file, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing file: Not a Golang source file?\n")
		os.Exit(-1)
	}

	// Inspect (DFS-walk) the AST
	// Anon. func is called on encounter of each node
	jv.Inspect(f, func(n ast.Node) bool {
		if debug {
			switch x := n.(type) {
			default:
				fmt.Printf("%T --> %+v\n", x, x)
			}
		}
		return true
	})
}
