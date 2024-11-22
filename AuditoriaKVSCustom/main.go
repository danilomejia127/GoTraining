package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type OSItem struct {
	Item struct {
		Key struct {
			S string `json:"S"`
		} `json:"key"`
		Metadata struct {
			S string `json:"S"`
		} `json:"metadata"`
		Version struct {
			N string `json:"N"`
		} `json:"version"`
		LastUpdated struct {
			S time.Time `json:"S"`
		} `json:"last_updated"`
		DateCreated struct {
			S time.Time `json:"S"`
		} `json:"date_created"`
		CompressedValue struct {
			B string `json:"B"`
		} `json:"compressed_value"`
		LastUpdatedMicros struct {
			N string `json:"N"`
		} `json:"last_updated_micros"`
	} `json:"Item"`
}

var mapData = make(map[string]bool)

func main() {
	// Populate gzipFiles with the names of the files in the 'data' folder
	files, err := ioutil.ReadDir("data")
	if err != nil {
		fmt.Println("Error reading the folder:", err)
		return
	}

	for _, file := range files {
		fmt.Println("Procesando file: ", file.Name())
		err := getDumpData("data/" + file.Name())
		if err != nil {
			fmt.Errorf("Error getting keys", err)

			return
		}
	}

	fmt.Println("Cantidad de keys: ", len(mapData))

	// Save keys to keysKVSCustom.txt
	err = saveKeysToFile("keysKVSCustom.txt", mapData)
	if err != nil {
		fmt.Println("Error saving keys to file:", err)
	}
}

func getDumpData(fileName string) error {
	// Abrir el archivo .gz
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error abriendo el archivo:", err)

		return err
	}
	defer file.Close()

	// Crear un lector gzip
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		fmt.Println("Error creando el lector gzip:", err)

		return err
	}
	defer gzReader.Close()

	scanner := bufio.NewScanner(gzReader)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		var item OSItem
		err := json.Unmarshal([]byte(scanner.Text()), &item)
		if err != nil {
			fmt.Errorf("failed Unmarshal", err)
			return err
		}

		/*decoded, err := base64.StdEncoding.DecodeString(item.Item.CompressedValue.B)
		if err != nil {
			fmt.Errorf("failed DecodeString", err)
			return nil, err
		}

		decompressed, err := gzip.NewReader(bytes.NewReader(decoded))
		if err != nil {
			fmt.Errorf("failed gzip.NewReader", err)
			return nil, err
		}

		value, err := ioutil.ReadAll(decompressed)
		if err != nil {
			return nil, err
		}*/

		// Imprimir el contenido
		// fmt.Println("Contenido del archivo descomprimido:", string(value))

		// fmt.Println("Key:", item.Item.Key.S)
		mapData[item.Item.Key.S] = true

	}

	return err
}

func saveKeysToFile(fileName string, data map[string]bool) error {
	file, err := os.Create("data/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for key := range data {
		_, err := writer.WriteString(key + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
