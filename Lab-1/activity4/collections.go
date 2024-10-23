package main

import "fmt"


func main(){
	fmt.Print("Name: Bisrat Asaye \nID: UGR/8508/14\n\n")
	
	// Array
	arr := [3]int{1,2,3}
	fmt.Println("Array: ", arr)

	// Slice
	slice := []int{4,5,6}
	slice = append(slice, 7)
	fmt.Println("Slice: ", slice)

	// Map(key-value pairs)
	myMap := map[string]int{} // make(map[string]int)
	myMap["Alice"] = 25
	myMap["Bob"] = 30
	fmt.Println("Map: ", myMap)
	fmt.Println("Alice's age: ", myMap["Alice"])

	// looping over a slice
	for i,v := range slice{
		fmt.Printf("index: %d, value %d\n", i,v)
	}

	// looping over a map
	for key,value := range myMap{
		fmt.Printf("%s is %d years old.\n", key, value)
	}
}

