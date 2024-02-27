package apicalls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/mercadolibre/GoTraining/compare_pir_vs_pmc/dtos"
)

const reportUsersCustom = "https://internal-api.mercadopago.com/v1/payment_methods/dump_users_custom"

var (
	Headers = map[string]string{}
)

func SetToken(t string) {
	Headers = map[string]string{
		"X-Tiger-Token": t,
	}
}

type RefreshMessage struct {
	Msg struct {
		ID struct {
			SiteID   string `json:"site_id"`
			SellerID string `json:"seller_id"`
		} `json:"id"`
	} `json:"msg"`
}

func CreateCustomData(refreshURL string, siteID string, sellerID int, scope string) bool {
	msg := RefreshMessage{}
	msg.Msg.ID.SiteID = siteID
	msg.Msg.ID.SellerID = strconv.Itoa(sellerID)

	jsonData, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return false
	}

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, refreshURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")

	for key, value := range Headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()

	// validar si el status code es 200
	if resp.StatusCode == 200 {
		return true
	}

	fmt.Println(fmt.Sprintf("Error al crear en scope: %s status: %d %s error: %s", scope, sellerID, resp.Status, resp.Body))

	return false
}

func GetDataCustomFromURL(url string, seller int) *dtos.Collector {
	urlFull := fmt.Sprintf("%s/v2/payment-methods/internal/customizations?collector_id=%d", url, seller)

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, urlFull, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for key, value := range Headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var dataCustom dtos.Collector

	err = json.Unmarshal(body, &dataCustom)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &dataCustom
}

func DeleteDataCustomFromURL(url string, seller int) bool {
	urlFull := url + fmt.Sprintf("%d", seller)

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, urlFull, nil)
	if err != nil {
		fmt.Println(err)
		return false
	}

	for key, value := range Headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return true
	}

	return false
}

func GetOriginalDataCustomFromURL(url string, seller int) *dtos.OriginalDataKVS {
	urlFull := fmt.Sprintf("%s%d", url, seller)

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, urlFull, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for key, value := range Headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var originalDataKVS dtos.OriginalDataKVS

	err = json.Unmarshal(body, &originalDataKVS)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &originalDataKVS
}

func GetSellerCustomFromReadV2() []int {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, reportUsersCustom, nil)
	if err != nil {
		fmt.Println(err)

		return nil
	}

	for key, value := range Headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)

		return nil
	}

	var sellers []int

	err = json.Unmarshal(body, &sellers)
	if err != nil {
		fmt.Println(err)

		return nil
	}

	return sellers
}
