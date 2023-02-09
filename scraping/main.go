package main

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

/* Para crear un programa en Go que explore las noticias de "www.elcolombiano.com" y guarde los titulares más relevantes,
podemos usar la librería "goquery" para hacer scraping de la página web y extraer la información que necesitamos.
*/

func main() {
	// URL de la página de noticias
	url := "https://www.elcolombiano.com"

	// Hacemos una petición HTTP para obtener la página web
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatalf("Error al obtener la página: %s", err)
	}

	// Un arreglo para guardar los titulares
	var titulares []string

	// Seleccionamos los elementos que contienen los titulares de noticias
	doc.Find("h3").Each(func(i int, s *goquery.Selection) {
		// Guardamos el texto del titular en el arreglo
		titulo := s.Text()
		titulares = append(titulares, titulo)
	})

	// Imprimimos los titulares en la consola
	fmt.Println("Titulares de noticias:")
	for _, titular := range titulares {
		fmt.Println(titular)
	}
}
