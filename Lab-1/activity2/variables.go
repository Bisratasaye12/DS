package main

import "fmt"

func main() {
	fmt.Print("Name: Bisrat Asaye \nID: UGR/8508/14\n\n")
	
	// Declaring variables with explicit types
	var a int = 10
	var b float64 = 20.5
	var c string = "Golang"
	var d bool = true

	// Short variable declaration (type inferred)
	e := "Short Declaration"

	// Print the values of the variables
	fmt.Println("Integer:", a)
	fmt.Println("Float:", b)
	fmt.Println("String:", c)
	fmt.Println("Boolean:", d)
	fmt.Println("Short Declaration:", e)
	
}
