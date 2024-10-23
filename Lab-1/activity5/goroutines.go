package main

import (
	"fmt"
	"time"
)


func PrintNumbers(){
	defer fmt.Print("Name: Bisrat Asaye \nID: UGR/8508/14\n\n")

	for i := 0; i <= 5; i++{
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
}

func PrintLetters(){
	for i := 'A'; i < 'E'; i++{
		fmt.Printf("%c\n", i)
		time.Sleep(1 * time.Second)
	}
}

func main(){

	go PrintNumbers()

	PrintLetters()

	time.Sleep(6 * time.Second)
	fmt.Println("Main Function Finished.")
}