package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	token      = ""
	tigerToken = "Bearer " + token
	apiUrl     = "https://testing-synchronizer-v2--payment-methods-read-v2.furyapps.io"
	resource   = "/pm-core/repository/custom-list"
)

func main() {
	start := time.Now()
	// Open the gzipped file
	gzippedFile, err := os.ReadFile("gzip_to_string/22wqaueduq3nhhagjvievoevim.json.gz")
	if err != nil {
		panic(err)
	}

	keyData, err := getKeyInfosFromFile(gzippedFile)
	if err != nil {
		fmt.Errorf("no es posible procesar %w", err)
		return
	}

	err = sendDataToKvs(keyData, 1000)
	if err != nil {
		return
	}

	fmt.Println("Fin " + strconv.Itoa(len(keyData)))
	end := time.Now()
	duration := end.Sub(start)
	fmt.Println("Tiempo transcurrido:", duration)
}

func sendDataToKvs(data []WrapperKvsV2, size int) error {
	for i := 0; i < len(data); i += size {
		end := i + size
		if end > len(data) {
			end = len(data)
		}
		subarray := data[i:end]

		fmt.Printf("Enviando registros del %v al %v \n", i, end)
		err := requestToReadV2(subarray)
		if err != nil {
			return err
		}
	}

	return nil
}

func requestToReadV2(data []WrapperKvsV2) error {
	requestBody, err := json.Marshal(data)
	if err != nil {
		return err
	}
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(requestBody))
	req.Header.Add("X-Tiger-Token", tigerToken)
	if err != nil {
		return err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("error con body %w \n", err)
	}

	fmt.Println(resp.Status)
	fmt.Println(string(bytes))

	return nil
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
			fmt.Errorf("failed Unmarshal %w", err)
			return []WrapperKvsV2{}, err
		}

		decoded, err := base64.StdEncoding.DecodeString(item.Item.CompressedValue.B)
		if err != nil {
			fmt.Errorf("failed DecodeString %w", err)
			return []WrapperKvsV2{}, err
		}

		decompressed, err := gzip.NewReader(bytes.NewReader(decoded))
		if err != nil {
			fmt.Errorf("failed gzip.NewReader %w", err)
			return []WrapperKvsV2{}, err
		}

		value, err := io.ReadAll(decompressed)
		if err != nil {
			return []WrapperKvsV2{}, err
		}

		var innerValue InnerValue

		err = json.Unmarshal(value, &innerValue)
		if err != nil {
			fmt.Errorf("failed json.Unmarshal InnerValue %w", err)
			return []WrapperKvsV2{}, err
		}

		keyInfo := WrapperKvsV2{
			Key:   item.Item.Key.S,
			Value: innerValue,
		}

		if keyInfo.Value.IsInStorage {
			fmt.Printf("In Storage %v\n", keyInfo.Key)
		} else {
			result = append(result, keyInfo)
		}
	}

	fmt.Println("getKeyInfosFromFile OK")

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
