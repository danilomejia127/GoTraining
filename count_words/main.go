package main

import (
	"fmt"
	"regexp"
	"strings"
)

/*
 * Reto #7
 * CONTANDO PALABRAS
 * Fecha publicación enunciado: 14/02/22
 * Fecha publicación resolución: 21/02/22
 * Dificultad: MEDIA
 *
 * Enunciado: Crea un programa que cuente cuantas veces se repite cada palabra y que muestre el recuento final de todas ellas.
 * - Los signos de puntuación no forman parte de la palabra.
 * - Una palabra es la misma aunque aparezca en mayúsculas y minúsculas.
 * - No se pueden utilizar funciones propias del lenguaje que lo resuelvan automáticamente.
 */

var nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9 ]+`)

func main() {
	input := "Betty bought the butter, the butter was bitter, " +
		"betty bought more butter to make the bitter butter better"
	// Todo a minusculas y dejar solo las letras y numeros
	input = clearString(strings.ToLower(input))

	fmt.Println("Input: " + input)
	for index, element := range repetition(input) {
		fmt.Println(index, "=", element)
	}
}

func repetition(str string) map[string]int {
	input := strings.Split(str, " ")
	wc := make(map[string]int)
	for _, word := range input {
		if _, matched := wc[word]; matched {
			wc[word] += 1
		} else {
			wc[word] = 1
		}
	}
	return wc
}

func clearString(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}
