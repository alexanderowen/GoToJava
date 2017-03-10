## GoToJava  
Translates a subset of the Go language to Java  
## Executable  
Compile and execute with   
`go build main.go`  
`./main filename`  
## Go Subset  
This program does not translate all Go programs.  
Current expectations of Go source  
- All variable declarations/definitions of the form `var x string = "hello"` or `var x string`  
- Go methods immediately follow the Go struct type definitions  
Example:  
`type Square struct {`  
`...`  
`}`  
`func (s Square) foo() {`  
`...`  
`}`    
- There is only one 'type struct' defintion per Go file. This conforms to Java syntax (not including anonymous classes)    
- Go does not allow initialization of the syntax `var i int = 0` in the init of a for loop delcaration  
As a result, for loops should be of the syntax  
`var i int`  
`for i = 0; i < n; i++ {`  
`...`  
`}`  

## Other Notes  
- Determines public/private based on Go capitalization syntax.  
- Determines if a Go function should map to a Java function based on if the Go package is 'main'.   
This behavior should be explored. May not preferable.   

