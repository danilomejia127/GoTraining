package main

import "fmt"

/*
 * Reto #3
 * ¿ES UN NÚMERO PRIMO?
 * Fecha publicación enunciado: 17/01/22
 * Fecha publicación resolución: 24/01/22
 * Dificultad: MEDIA
 *
 * Enunciado: Escribe un programa que se encargue de comprobar si un número es o no primo. Los primos son los que solo se dejan dividir por 1 y por si mismo
 * Hecho esto, imprime los números primos entre 1 y 100.
 *
 */

func main() {
	for i := 1; i <= 100; i++ {
		if isPrime(i) {
			fmt.Println(i)
		}
	}
}

func isPrime(number int) bool {
	if number < 2 {
		return true
	}
	for i := 2; i < number; i++ {
		if number%i == 0 {
			return false
		}
	}

	return true
}
