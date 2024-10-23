package main

import "fmt"

func add(x int, y int) int {
	return x + y;
}

func main(){
	fmt.Print("Name: Bisrat Asaye \nID: UGR/8508/14\n\n")
	
	result := add(5, 3)

	fmt.Println("sum: ", result)


	if result > 5 {
		fmt.Println("The result is greater than 5.")
	} else if result == 5{
		fmt.Println("The result is equal to 5.")
	} else{
		fmt.Println("The result is less than 5.")
	}


	for i := 0; i < 5; i++{
		fmt.Println("Loop iteration: ", i)
	}
}