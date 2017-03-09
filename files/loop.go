package main

import "fmt"

func main() {
	var x int = 0
	for x < 10 {
		x = x + 1
	}
	var i int = 0
	for ; i < 10; i += 1 {
		fmt.Println(i)
	}
	for {
		fmt.Println("Infinite loop")
	}
}
