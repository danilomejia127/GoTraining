package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Key struct {
	S string `json:"S"`
}

type Metadata struct {
	S string `json:"S"`
}

type Version struct {
	N string `json:"N"`
}

type LastUpdated struct {
	S string `json:"S"`
}

type DateCreated struct {
	S string `json:"S"`
}

type CompressedValue struct {
	B string `json:"B"`
}

type LastUpdatedMicros struct {
	N string `json:"N"`
}

type Item struct {
	Key               Key               `json:"key"`
	Metadata          Metadata          `json:"metadata"`
	Version           Version           `json:"version"`
	LastUpdated       LastUpdated       `json:"last_updated"`
	DateCreated       DateCreated       `json:"date_created"`
	CompressedValue   CompressedValue   `json:"compressed_value"`
	LastUpdatedMicros LastUpdatedMicros `json:"last_updated_micros"`
}

type Data struct {
	Item Item `json:"Item"`
}

func main() {
	start := time.Now()

	file, err := os.Open("ReadObjectStorageMeliFile/7ojvr3ozl43qvdkkpu4ujjwqzu.json")
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer file.Close()

	db, err := sql.Open("mysql", "root:pmlocal@tcp(localhost:3306)/pmdev")
	if err != nil {
		fmt.Println("Error al conectar a la base de datos: ", err)
		return
	}
	defer db.Close()

	// Verificar la conexión
	err = db.Ping()
	if err != nil {
		fmt.Println("Error al establecer la conexión: ", err)
		return
	}

	dataCh := make(chan Data)
	progressCh := make(chan bool)

	var wg sync.WaitGroup
	var progressWg sync.WaitGroup

	// Crear 20 workers
	for i := 0; i < 20; i++ {
		wg.Add(1)
		progressWg.Add(1)
		go worker(db, dataCh, progressCh, &wg, &progressWg)
	}

	// Goroutine para rastrear el progreso
	go func() {
		processed := 0
		for range progressCh {
			processed++
			fmt.Printf("Processed %d records\n", processed)
		}
	}()

	scanner := bufio.NewScanner(file)
	processed := 0

	for scanner.Scan() {
		var data Data
		err := json.Unmarshal(scanner.Bytes(), &data)
		if err != nil {
			fmt.Println("Error decoding line: ", err)
			continue
		}

		// Enviar los datos a las goroutines
		dataCh <- data
		processed++
	}

	close(dataCh)
	wg.Wait()
	close(progressCh)

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}

	elapsed := time.Since(start)
	minutes := int(elapsed.Minutes())
	seconds := int(elapsed.Seconds()) % 60
	fmt.Printf("Processed %d records in %d minutes and %d seconds\n", processed, minutes, seconds)
}

func worker(db *sql.DB, dataCh <-chan Data, progressCh chan<- bool, wg *sync.WaitGroup, progressWg *sync.WaitGroup) {
	defer wg.Done()

	// Preparar la consulta SQL
	stmt, err := db.Prepare("INSERT INTO custom_seller(id, site_id, in_prod, in_clon) VALUES(?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Error al preparar la consulta: ", err)
		return
	}
	defer stmt.Close()

	for data := range dataCh {
		inProd := 0
		inClon := 0
		_, err = stmt.Exec(data.Item.Key.S, "", inProd, inClon)
		if err != nil {
			fmt.Println("Error al insertar el valor en la base de datos: ", err)
			continue
		}

		// Notificar que se ha procesado un registro
		progressCh <- true
	}

	progressWg.Done()
}
