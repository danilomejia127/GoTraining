package main

import "fmt"

/*
 * Reto #4
 * ÁREA DE UN POLÍGONO
 * Fecha publicación enunciado: 24/01/22
 * Fecha publicación resolución: 31/01/22
 * Dificultad: FÁCIL
 *
 * Enunciado: Crea UNA ÚNICA FUNCIÓN (importante que sólo sea una) que sea capaz de calcular y retornar el área de un polígono.
 * - La función recibirá por parámetro sólo UN polígono a la vez.
 * - Los polígonos soportados serán Triángulo, Cuadrado y Rectángulo.
 * - Imprime el cálculo del área de un polígono de cada tipo.
 *
 */

/* Interface dacaration type Polygon */
type Polygon interface {
	area() float64
}

type Rectangle struct {
	base, height float64
}

func (polygon *Rectangle) area() float64 {
	return polygon.base * polygon.height
}

type Square struct {
	side float64
}

func (polygon *Square) area() float64 {
	return polygon.side * polygon.side
}

type Triangle struct {
	base, height float64
}

func (polygon *Triangle) area() float64 {
	return (polygon.base * polygon.height) / 2
}

/* This is the function that can calcualte area from any poligon defined*/
func gimeArea(fig Polygon) float64 {
	return fig.area()
}

func main() {
	rectangle := Rectangle{2.1, 4.3}
	square := Square{2.3}
	triangle := Triangle{5.3, 2.8}

	fmt.Println("Rectangle area is: ", gimeArea(&rectangle))
	fmt.Println("Square area is: ", gimeArea(&square))
	fmt.Println("Triangle area is: ", gimeArea(&triangle))
}
