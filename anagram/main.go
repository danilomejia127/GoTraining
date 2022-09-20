package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
 * Reto #1
 * ¿ES UN ANAGRAMA?
 * Fecha publicación enunciado: 03/01/22
 * Fecha publicación resolución: 10/01/22
 * Dificultad: MEDIA
 *
 * Enunciado: Escribe una función que reciba dos palabras (String) y retorne verdadero o falso (Boolean) según sean o no anagramas.
 * Un Anagrama consiste en formar una palabra reordenando TODAS las letras de otra palabra inicial.
 * NO hace falta comprobar que ambas palabras existan.
 * Dos palabras exactamente iguales no son anagrama.
 *
 */

func main() {
	var string1, string2, sortedString1, sortedString2 string

	fmt.Println("Enter 1st word")
	fmt.Scanf("%s", &string1)

	fmt.Println("Enter 2nd word")
	fmt.Scanf("%s", &string2)

	string1Slice := strings.Split(string1, "")
	string2Slice := strings.Split(string2, "")

	sort.Slice(string1Slice, func(i, j int) bool {
		return string1Slice[i] < string1Slice[j]
	})

	sort.Slice(string2Slice, func(i, j int) bool {
		return string2Slice[i] < string2Slice[j]
	})

	sortedString1 = strings.Join(string1Slice, "")
	sortedString2 = strings.Join(string2Slice, "")

	if sortedString1 == sortedString2 {
		fmt.Println("These words are anagrams")
	} else {
		fmt.Println("These words are not anagrams")
	}
}
