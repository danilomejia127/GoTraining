package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	// Open the gzipped file
	gzippedFile, err := os.ReadFile("gzip_to_string/olvue4mdwa32jcfi24jkpq7nvy.json.gz")
	if err != nil {
		panic(err)
	}

	keyData, err := getKeyInfosFromFile(gzippedFile)

	fmt.Println("Fin " + strconv.Itoa(len(keyData)))

}

func getKeyInfosFromFile(raw []byte) ([]KeyInfo, error) {
	gzreader, err := gzip.NewReader(bytes.NewReader(raw))
	if err != nil {
		fmt.Printf("failed gzip.NewReader %v", err)
		return []KeyInfo{}, err
	}

	scanner := bufio.NewScanner(gzreader)
	scanner.Split(bufio.ScanLines)

	result := make([]KeyInfo, 0, 10)

	for scanner.Scan() {
		var item OSItem
		err := json.Unmarshal([]byte(scanner.Text()), &item)
		if err != nil {
			fmt.Printf("failed Unmarshal %v", err)
			return []KeyInfo{}, err
		}

		decoded, err := base64.StdEncoding.DecodeString(item.Item.CompressedValue.B)
		if err != nil {
			fmt.Printf("failed DecodeString %v", err)
			return []KeyInfo{}, err
		}

		decompressed, err := gzip.NewReader(bytes.NewReader(decoded))
		if err != nil {
			fmt.Printf("failed gzip.NewReader %v", err)
			return []KeyInfo{}, err
		}

		value, err := io.ReadAll(decompressed)
		if err != nil {
			return []KeyInfo{}, err
		}

		keyInfo := KeyInfo{
			Key:  item.Item.Key.S,
			Info: string(value),
		}

		result = append(result, keyInfo)
	}

	fmt.Println("getKeyInfosFromFile succeeded")

	return result, nil
}

type OSItem struct {
	Item struct {
		CompressedValue struct {
			B string `json:"B"`
		} `json:"compressed_value"`
		Key struct {
			S string `json:"S"`
		} `json:"key"`
	} `json:"Item"`
}

type KeyInfo struct {
	Key  string
	Info string
}
