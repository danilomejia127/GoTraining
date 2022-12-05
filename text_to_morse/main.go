package main

import (
	"fmt"
)

/*
 * Reto #9
 * CÓDIGO MORSE
 * Fecha publicación enunciado: 02/03/22
 * Fecha publicación resolución: 07/03/22
 * Dificultad: MEDIA
 *
 * Enunciado: Crea un programa que sea capaz de transformar texto natural a código morse y viceversa.
 * - Debe detectar automáticamente de qué tipo se trata y realizar la conversión.
 * - En morse se soporta raya "—", punto ".", un espacio " " entre letras o símbolos y dos espacios entre palabras "  ".
 * - El alfabeto morse soportado será el mostrado en https://es.wikipedia.org/wiki/Código_morse.
 *
 */
func main() {

	for k, elem := range getMorseCode() {
		fmt.Println(k + "=" + elem)
	}

}

func getMorseCode() map[string]string {
	var codeMorse = make(map[string]string)
	codeMorse["A"] = "-."
	codeMorse["B"] = "-..."
	codeMorse["C"] = "-.-."
	codeMorse["D"] = "-.."
	codeMorse["E"] = "."
	codeMorse["F"] = "..-."
	codeMorse["G"] = "--."
	codeMorse["H"] = "...."
	codeMorse["I"] = "..-"
	codeMorse["J"] = ".---"
	codeMorse["K"] = "-.."
	codeMorse["L"] = ".-.."
	codeMorse["M"] = "--"
	codeMorse["Ñ"] = "--.--"
	codeMorse["N"] = "-."
	codeMorse["O"] = "---"
	codeMorse["P"] = ".--."
	codeMorse["Q"] = "--.-"
	codeMorse["R"] = ".-."
	codeMorse["S"] = "..."
	codeMorse["T"] = "-"
	codeMorse["U"] = ".."
	codeMorse["V"] = "...-"
	codeMorse["W"] = ".--"
	codeMorse["X"] = "-..-"
	codeMorse["Y"] = "-.--"
	codeMorse["Z"] = "--.."
	codeMorse["0"] = "------"
	codeMorse["1"] = ".----"
	codeMorse["2"] = "..---"
	codeMorse["3"] = "...--"
	codeMorse["4"] = "....-"
	codeMorse["5"] = "....."
	codeMorse["6"] = "-...."
	codeMorse["7"] = "--..."
	codeMorse["8"] = "---.."
	codeMorse["9"] = "----."

	return codeMorse
}
