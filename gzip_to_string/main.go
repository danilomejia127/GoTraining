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
	"time"
)

func main() {
	// Open the gzipped file
	gzippedFile, err := os.ReadFile("gzip_to_string/22wqaueduq3nhhagjvievoevim.json.gz")
	if err != nil {
		panic(err)
	}

	keyData, err := getKeyInfosFromFile(gzippedFile)

	fmt.Println("Fin " + strconv.Itoa(len(keyData)))

}

func getKeyInfosFromFile(raw []byte) ([]WrapperKvsV2, error) {
	gzreader, err := gzip.NewReader(bytes.NewReader(raw))
	if err != nil {
		fmt.Printf("failed gzip.NewReader %v", err)
		return []WrapperKvsV2{}, err
	}

	scanner := bufio.NewScanner(gzreader)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	scanner.Split(bufio.ScanLines)

	var result []WrapperKvsV2

	for scanner.Scan() {
		var item OSItem
		err := json.Unmarshal([]byte(scanner.Text()), &item)
		if err != nil {
			fmt.Printf("failed Unmarshal %v", err)
			return []WrapperKvsV2{}, err
		}

		decoded, err := base64.StdEncoding.DecodeString(item.Item.CompressedValue.B)
		if err != nil {
			fmt.Printf("failed DecodeString %v", err)
			return []WrapperKvsV2{}, err
		}

		decompressed, err := gzip.NewReader(bytes.NewReader(decoded))
		if err != nil {
			fmt.Printf("failed gzip.NewReader %v", err)
			return []WrapperKvsV2{}, err
		}

		value, err := io.ReadAll(decompressed)
		if err != nil {
			return []WrapperKvsV2{}, err
		}

		var innerValue InnerValue

		err = json.Unmarshal(value, &innerValue)
		if err != nil {
			fmt.Printf("failed json.Unmarshal InnerValue %v", err)
			return []WrapperKvsV2{}, err
		}

		keyInfo := WrapperKvsV2{
			Key:   item.Item.Key.S,
			Value: innerValue,
		}

		if keyInfo.Value.IsInStorage {
			fmt.Printf("In Storage %v\n", keyInfo.Key)
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

type WrapperKvsV2 struct {
	Key   string
	Value InnerValue
}

type InnerValue struct {
	IsCompressed bool
	IsInStorage  bool
	Data         string
	LastUpdated  time.Time
}
