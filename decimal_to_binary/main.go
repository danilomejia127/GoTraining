package main

import "fmt"

func main() {
	decimalToBinary(1985)
	decimalToBinary(2021)
	decimalToBinary(10)
	decimalToBinary(31)
	decimalToBinary(0)
}

func decimalToBinary(decimal int) {
	var binary []int
	for decimal != 0 {
		binary = append(binary, decimal%2)
		// decimal = decimal / 2
		decimal /= 2
	}

	if len(binary) == 0 {
		fmt.Printf("%d\n", 0)
	} else {
		for i := len(binary) - 1; i >= 0; i-- {
			fmt.Printf("%d", binary[i])
		}
		fmt.Println()
	}
}
