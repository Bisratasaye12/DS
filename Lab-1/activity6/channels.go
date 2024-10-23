package main

import "fmt"

func SendData(ch chan int){
	for i := 0; i < 5; i++{
		ch <- i
	}
	close(ch)
}

func main(){
	fmt.Print("Name: Bisrat Asaye \nID: UGR/8508/14\n\n")

	ch := make(chan int)
	go SendData(ch)

	for val := range ch{
		fmt.Println("Recieved: ", val)
	}

	fmt.Println("Channel Closed, Program Finished!")
}