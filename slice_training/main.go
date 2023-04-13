package main

import "fmt"

func main() {
	numbers := []int{2, 3, 4, 5, 6, 8, 10, 20, 2}
	fmt.Printf("The small number is: %v", smallestNumber(numbers))
}
func smallestNumber(slice []int) int {
	small := slice[0]
	for _, number := range slice {
		if number < small {
			small = number
		}
	}
	return small
}
