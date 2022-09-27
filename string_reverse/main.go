package main

import "fmt"

/*
 * Reto #6
 * INVIRTIENDO CADENAS
 * Fecha publicación enunciado: 07/02/22
 * Fecha publicación resolución: 14/02/22
 * Dificultad: FÁCIL
 *
 * Enunciado: Crea un programa que invierta el orden de una cadena de texto sin usar funciones propias del lenguaje que lo hagan de forma automática.
 * - Si le pasamos "Hola mundo" nos retornaría "odnum aloH"
 */
func main() {
	fmt.Println(reverseString("INVIRTIENDO CADENAS"))
	fmt.Println(reverseString("Danilo"))
}

func reverseString(str string) string {
	byteStr := []rune(str)
	for i, j := 0, len(byteStr)-1; i < j; i, j = i+1, j-1 {
		byteStr[i], byteStr[j] = byteStr[j], byteStr[i]
	}
	return string(byteStr)
}
