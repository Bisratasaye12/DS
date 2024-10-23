package main

import (
  "fmt"
  "os"
)


// Function that opens a file and reads it (simplified for this example)
func readFile(filename string) error {
  file, err := os.Open(filename)
  if err != nil {
    return err // Return the error if file opening fails
  }
  defer file.Close() // Ensure the file is closed after function returns
  fmt.Println("File opened successfully:", filename)
  return nil
}
func main() {
  fmt.Print("Name: Bisrat Asaye \nID: UGR/8508/14\n\n")
  
	err := readFile("test.txt")
	if err != nil {
	fmt.Println("Error:", err)
	} else {
	fmt.Println("File read successfully.")
	}
}