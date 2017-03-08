/* Custom visitor for the Go compiler AST
 * For a subset of possible Go programs, prints an approximation of the source to stdout
 * Large part of code base provided by the Go compiler visitor, "go/ast/walk.go"
 * */

package JavaVisitor

import (
	"fmt"
	"go/ast"
	"strings"
)

// A Visitor's Visit method is invoked for each node encountered by Walk.
// If the result visitor w is not nil, Walk visits each of the children
// of node with the visitor w, followed by a call of w.Visit(nil).
type Visitor interface {
	Visit(node ast.Node) (w Visitor)
}

// javaMap specifies certain 1-1 relationships between go and java
// libraries, functions, keywords, and so on
// May not necessarily hold for all cases
var javaMap = map[string]string{
	"fmt":     "System.out",
	"Println": "println",
	//	"func":    "public static",
	"type":   "class",
	"string": "String",
}

var defMethods bool = false
var curCls string

// Helper functions for common node lists. They may be empty.
func walkIdentList(v Visitor, list []*ast.Ident) {
	for _, x := range list {
		Walk(v, x)
	}
}

func walkExprList(v Visitor, list []ast.Expr) {
	for i, x := range list {
		Walk(v, x)
		if i < len(list)-1 {
			fmt.Printf(", ")
		}
	}
}

func walkStmtList(v Visitor, list []ast.Stmt) {
	for _, x := range list {
		Walk(v, x)
		/*
			if i < len(list)-1 {
				fmt.Printf(";\n")
			}*/
		fmt.Printf(";\n")
	}
}

func walkDeclList(v Visitor, list []ast.Decl) {
	for _, x := range list {
		Walk(v, x)
	}
}

// Walk traverses an AST in depth-first order: It starts by calling
// v.Visit(node); node must not be nil. If the visitor w returned by
// v.Visit(node) is not nil, Walk is invoked recursively with visitor
// w for each of the non-nil children of node, followed by a call of
// w.Visit(nil).
//
func Walk(v Visitor, node ast.Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	// walk children
	// (the order of the cases matches the order
	// of the corresponding node types in ast.go)
	switch n := node.(type) {
	// Comments and fields
	case *ast.Comment:
		// nothing to do

	case *ast.CommentGroup:
		for _, c := range n.List {
			Walk(v, c)
		}

	case *ast.Field:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		Walk(v, n.Type) //Java, type first
		fmt.Printf(" ")
		walkIdentList(v, n.Names)
		if n.Tag != nil {
			Walk(v, n.Tag)
		}
		if n.Comment != nil {
			Walk(v, n.Comment)
		}
		fmt.Printf(";")

	case *ast.FieldList:
		//fmt.Printf("(")
		for _, f := range n.List {
			Walk(v, f)
			/*
				if i != len(n.List)-1 {
					fmt.Printf(", ")
				}
			*/
		}
		//fmt.Printf(")")

	// Expressions
	case *ast.BadExpr:
		// do nothing

	case *ast.Ident:
		if val, ok := javaMap[n.Name]; ok {
			fmt.Printf("%s", val)
		} else {
			fmt.Printf("%s", n.Name)
		}

	case *ast.BasicLit:
		fmt.Printf("%s", n.Value)

	case *ast.Ellipsis:
		if n.Elt != nil {
			Walk(v, n.Elt)
		}

	case *ast.FuncLit:
		fmt.Printf("func")
		Walk(v, n.Type)
		fmt.Printf(" {\n")
		Walk(v, n.Body)
		fmt.Printf("\n}")

	case *ast.CompositeLit:
		if n.Type != nil {
			Walk(v, n.Type)
		}
		fmt.Printf("{")
		walkExprList(v, n.Elts)
		fmt.Printf("}")

	case *ast.ParenExpr:
		Walk(v, n.X)

	case *ast.SelectorExpr:
		fmt.Printf(" ")
		Walk(v, n.X)
		fmt.Printf(".")
		Walk(v, n.Sel)

	case *ast.IndexExpr:
		Walk(v, n.X)
		fmt.Printf("[")
		Walk(v, n.Index)
		fmt.Printf("]")

	case *ast.SliceExpr:
		Walk(v, n.X)
		if n.Low != nil {
			Walk(v, n.Low)
		}
		if n.High != nil {
			Walk(v, n.High)
		}
		if n.Max != nil {
			Walk(v, n.Max)
		}

	case *ast.TypeAssertExpr:
		Walk(v, n.X)
		fmt.Printf(".(")
		if n.Type != nil {
			Walk(v, n.Type)
		} else {
			fmt.Printf("type")
		}
		fmt.Printf(")")

	case *ast.CallExpr:
		Walk(v, n.Fun)
		fmt.Printf("(")
		walkExprList(v, n.Args)
		fmt.Printf(")")

	case *ast.StarExpr:
		Walk(v, n.X)

	case *ast.UnaryExpr:
		Walk(v, n.X)

	case *ast.BinaryExpr:
		Walk(v, n.X)
		fmt.Printf(" %s ", n.Op.String())
		Walk(v, n.Y)

	case *ast.KeyValueExpr:
		Walk(v, n.Key)
		Walk(v, n.Value)

	// Types
	case *ast.ArrayType:
		if n.Len != nil {
			fmt.Printf("[")
			Walk(v, n.Len)
			fmt.Printf("]")
		}
		Walk(v, n.Elt)

	case *ast.StructType:
		fmt.Printf(" {\n")
		Walk(v, n.Fields)
		//fmt.Printf("\n}")

	case *ast.FuncType:
		/* Handled in FuncDecl, so Java function ordering is correct
		if n.Results != nil {
			Walk(v, n.Results)
		} else {
			fmt.Printf(" void ")
		}
		*/
		if n.Params != nil {
			Walk(v, n.Params)
		}

	case *ast.InterfaceType:
		Walk(v, n.Methods)

	case *ast.MapType:
		Walk(v, n.Key)
		Walk(v, n.Value)

	case *ast.ChanType:
		Walk(v, n.Value)

	// Statements
	case *ast.BadStmt:
		// nothing to do

	case *ast.DeclStmt:
		Walk(v, n.Decl)

	case *ast.EmptyStmt:
		// nothing to do

	case *ast.LabeledStmt:
		Walk(v, n.Label)
		Walk(v, n.Stmt)

	case *ast.ExprStmt:
		Walk(v, n.X)

	case *ast.SendStmt:
		Walk(v, n.Chan)
		Walk(v, n.Value)

	case *ast.IncDecStmt:
		Walk(v, n.X)

	case *ast.AssignStmt:
		walkExprList(v, n.Lhs)
		fmt.Printf(" = ")
		walkExprList(v, n.Rhs)

	case *ast.GoStmt:
		Walk(v, n.Call)

	case *ast.DeferStmt:
		Walk(v, n.Call)

	case *ast.ReturnStmt:
		fmt.Printf("return ")
		walkExprList(v, n.Results)

	case *ast.BranchStmt:
		if n.Label != nil {
			Walk(v, n.Label)
		}

	case *ast.BlockStmt:
		walkStmtList(v, n.List)

	case *ast.IfStmt:
		fmt.Printf("if ")
		if n.Init != nil {
			Walk(v, n.Init)
		}
		Walk(v, n.Cond)
		fmt.Printf("{\n")
		Walk(v, n.Body)
		fmt.Printf("\n}")
		if n.Else != nil {
			fmt.Printf("else {\n")
			Walk(v, n.Else)
			fmt.Printf("}\n")
		}

	case *ast.CaseClause:
		if n.List == nil {
			fmt.Printf("default")
		} else {
			fmt.Printf("case ")
		}
		walkExprList(v, n.List)
		fmt.Printf(":\n")
		walkStmtList(v, n.Body)

	case *ast.SwitchStmt:
		fmt.Printf("switch ")
		if n.Init != nil {
			Walk(v, n.Init)
			fmt.Printf(";")
		}
		if n.Tag != nil {
			Walk(v, n.Tag)
		}
		fmt.Printf(" {\n")
		Walk(v, n.Body)
		fmt.Printf("\n}")

	case *ast.TypeSwitchStmt:
		fmt.Printf("switch ")
		if n.Init != nil {
			Walk(v, n.Init)
			fmt.Printf(";")
		}
		Walk(v, n.Assign)
		fmt.Printf(" {\n")
		Walk(v, n.Body)
		fmt.Printf("\n}")

	case *ast.CommClause:
		if n.Comm != nil {
			Walk(v, n.Comm)
		}
		walkStmtList(v, n.Body)

	case *ast.SelectStmt:
		Walk(v, n.Body)

	case *ast.ForStmt:
		if n.Init != nil {
			Walk(v, n.Init)
		}
		if n.Cond != nil {
			Walk(v, n.Cond)
		}
		if n.Post != nil {
			Walk(v, n.Post)
		}
		Walk(v, n.Body)

	case *ast.RangeStmt:
		fmt.Printf("for ")
		if n.Key != nil {
			Walk(v, n.Key)
		}
		if n.Value != nil {
			fmt.Printf(", ")
			Walk(v, n.Value)
		}
		fmt.Printf(":= range ")
		Walk(v, n.X)
		fmt.Printf(" {\n")
		Walk(v, n.Body)
		fmt.Printf("\n}")

	// Declarations
	case *ast.ImportSpec:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		if n.Name != nil {
			Walk(v, n.Name)
			fmt.Printf(" ")
		}
		Walk(v, n.Path)
		if n.Comment != nil {
			Walk(v, n.Comment)
		}

	case *ast.ValueSpec:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		if n.Type != nil { //walk the type first, for Java
			Walk(v, n.Type)
		}
		fmt.Printf(" ")
		walkIdentList(v, n.Names)
		if n.Values != nil {
			fmt.Printf(" = ")
		}
		walkExprList(v, n.Values)
		if n.Comment != nil {
			Walk(v, n.Comment)
		}
		//fmt.Printf("\n")

	case *ast.TypeSpec:
		defMethods = true
		curCls = n.Name.Name
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		fmt.Printf("%s ", javaMap["type"])
		Walk(v, n.Name)
		Walk(v, n.Type)
		if n.Comment != nil {
			Walk(v, n.Comment)
		}

	case *ast.BadDecl:
		// nothing to do

	case *ast.GenDecl:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}

		if n.Tok.String() == "import" { // ignore imports for now
			return
		} else if n.Tok.String() == "var" {
			//fmt.Printf("%s ", n.Specs[0].Values[0].Kind)
		}
		/*
			fmt.Printf("%s ", n.Tok)
			multiImport := n.Tok.String() == "import" && n.Lparen.IsValid()
			if multiImport {
				fmt.Printf("(")
			}
			for i, s := range n.Specs {
				Walk(v, s)
					if multiImport && i != len(n.Specs)-1 {
						fmt.Printf("\n")
					}
			}
				if multiImport {
					fmt.Printf(")\n")
				}
		*/
		for _, s := range n.Specs {
			Walk(v, s)
		}

	case *ast.FuncDecl:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		if n.Recv != nil {
			var t string = fmt.Sprintf("%v", (n.Recv.List[0].Type))
			if defMethods && curCls != t { //closes the Java class defintion
				defMethods = false
				curCls = ""
				fmt.Printf("}")
			}
			javaMap[n.Recv.List[0].Names[0].Name] = "this"
			/*
				fmt.Printf("(")
				Walk(v, n.Recv)
				fmt.Printf(")")
			*/
		}
		if n.Name.Name == "main" { //special Java function signature
			fmt.Printf("public static void main(String[] args) ")
		} else {
			//fmt.Printf("%s ", javaMap["func"])
			if string(n.Name.Name[0]) == strings.ToUpper(string(n.Name.Name[0])) {
				fmt.Printf("public ")
			} else {
				fmt.Printf("private ")
			}
			if n.Type.Results == nil {
				fmt.Printf("void ")
			} else {
				retVal := fmt.Sprintf("%s", n.Type.Results.List[0].Type)
				if val, ok := javaMap[retVal]; ok {
					fmt.Printf(" %s ", val)
				} else {
					fmt.Printf(" %s ", retVal)
				}
			}
			Walk(v, n.Name)
			fmt.Printf("(")
			Walk(v, n.Type)
			fmt.Printf(")")
		}
		fmt.Printf(" {\n")
		if n.Body != nil {
			Walk(v, n.Body)
		}
		fmt.Printf("\n}")

	// Files and packages
	case *ast.File:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}

		if n.Name.Name == "main" { //if Go package is main, map to 'class Main'
			fmt.Printf("class ")
			Walk(v, n.Name)
			fmt.Printf("{\n")
			walkDeclList(v, n.Decls)
			fmt.Printf("\n}")
		} else { //otherwise, treat as package
			fmt.Printf("package ")
			Walk(v, n.Name)
			fmt.Printf(";\n")
			walkDeclList(v, n.Decls)
			if defMethods { //closes Java class definition
				fmt.Printf("}")
			}
		}

	case *ast.Package:
		for _, f := range n.Files {
			Walk(v, f)
		}

	default:
		panic(fmt.Sprintf("ast.Walk: unexpected node type %T", n))
	}

	v.Visit(nil)
}

type inspector func(ast.Node) bool

func (f inspector) Visit(node ast.Node) Visitor {
	if f(node) {
		return f
	}
	return nil
}

// Inspect traverses an AST in depth-first order: It starts by calling
// f(node); node must not be nil. If f returns true, Inspect invokes f
// recursively for each of the non-nil children of node, followed by a
// call of f(nil).
//
func Inspect(node ast.Node, f func(ast.Node) bool) {
	Walk(inspector(f), node)
}
