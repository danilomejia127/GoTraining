package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	url1 := "URI1"
	url2 := "URI2"

	headers := map[string]string{
		"X-Tiger-Token": "",
	}

	result1, err := getRequest(url1, headers)
	if err != nil {
		fmt.Printf("Error al obtener los datos de URL1: %s\n", err)
		return
	}

	result2, err := getRequest(url2, headers)
	if err != nil {
		fmt.Printf("Error al obtener los datos de URL2: %s\n", err)
		return
	}

	difference := compareArrays(result1, result2)

	// Guardar la diferencia en un archivo
	err = saveToFile("diferencia.txt", difference)
	if err != nil {
		fmt.Printf("Error al guardar la diferencia en el archivo: %s\n", err)
		return
	}

	fmt.Println("Diferencia entre URL1 y URL2 guardada en diferencia.txt")
}

func getRequest(url string, headers map[string]string) ([]int, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result []int
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func compareArrays(arr1, arr2 []int) []int {
	difference := make([]int, 0)

	map1 := make(map[int]bool)
	for _, num := range arr1 {
		map1[num] = true
	}

	map2 := make(map[int]bool)
	for _, num := range arr2 {
		map2[num] = true
	}

	fmt.Printf("map1: %d, map2: %d\n", len(map1), len(map2))

	i := 0
	for num := range map2 {

		if !map1[num] {
			difference = append(difference, num)
			fmt.Printf("No existe : %d\n", num)
		}

		if i%100 == 0 {
			fmt.Printf("Revisando : %d\n", i)
		}

		i++
	}

	return difference
}

func saveToFile(filename string, data []int) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, value := range data {
		_, err := file.WriteString(fmt.Sprintf("%d\n", value))
		if err != nil {
			return err
		}
	}

	return nil
}
