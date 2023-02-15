package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1) // Añadir una goroutine al WaitGroup
		go func(id int) {
			defer wg.Done()                                       // Indicar que la goroutine ha finalizado
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second) // Esperar un tiempo aleatorio
			fmt.Printf("Goroutine %d ha terminado\n", id)
		}(i)
	}
	wg.Wait() // Esperar a que todas las goroutines finalicen
	fmt.Println("Todas las goroutines han terminado")
}
