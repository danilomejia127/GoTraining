package apicalls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/mercadolibre/GoTraining/compare_pir_vs_pmc/dtos"
)

const (
	reportUsersCustom            = "https://internal-api.mercadopago.com/v1/payment_methods/dump_users_custom"
	productionSyncReaderV2       = "https://production-synchronizer-v2--payment-methods-read-v2.furyapps.io"
	specialOwnersPath            = "/pm-core/repository/get-special-owners/"
	getSpecialOwnersBySiteURL    = productionSyncReaderV2 + specialOwnersPath
	ProductionSynchronizerStgURL = "https://production-synchronizer-stg--payment-methods-read-v2.furyapps.io/pm-core/repository/standard"
)

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

type SpecialOwnersMsg struct {
	Msg struct {
		Key   string `json:"Key"`
		Value []struct {
			ID     int      `json:"id"`
			Values []string `json:"values"`
		} `json:"Value"`
		MessageID   string `json:"message_id"`
		PublishTime int64  `json:"publish_time"`
	} `json:"msg"`
}

func CreateCustomData(refreshURL string, siteID string, sellerID int, scope string) bool {
	msg := RefreshMessage{}
	msg.Msg.ID.SiteID = siteID
	msg.Msg.ID.SellerID = strconv.Itoa(sellerID)

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)

		return false
	}

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, refreshURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)

		return false
	}

	req.Header.Set("Content-Type", "application/json")

	for key, value := range Headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)

		return false
	}
	defer resp.Body.Close()

	// validar si el status code es 200
	if resp.StatusCode == 200 {
		return true
	}

	log.Println(fmt.Sprintf("Error al crear en scope: %s status: %d %s error: %s", scope, sellerID, resp.Status, resp.Body))

	return false
}

func GetDataCustomFromURL(url string, seller int) *dtos.Collector {
	urlFull := fmt.Sprintf("%s/v2/payment-methods/internal/customizations?collector_id=%d", url, seller)

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, urlFull, nil)
	if err != nil {
		log.Println(err)

		return nil
	}

	for key, value := range Headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)

		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)

		return nil
	}

	var dataCustom dtos.Collector

	err = json.Unmarshal(body, &dataCustom)
	if err != nil {
		log.Println(err)

		return nil
	}

	return &dataCustom
}

func DeleteDataCustomFromURL(url string, seller int) bool {
	urlFull := url + fmt.Sprintf("%d", seller)

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, urlFull, nil)
	if err != nil {
		log.Println(err)

		return false
	}

	for key, value := range Headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)

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
		log.Println(err)

		return nil
	}

	for key, value := range Headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)

		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)

		return nil
	}

	var originalDataKVS dtos.OriginalDataKVS

	err = json.Unmarshal(body, &originalDataKVS)
	if err != nil {
		log.Println(err)

		return nil
	}

	return &originalDataKVS
}

func GetSellerCustomFromReadV2() []int {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, reportUsersCustom, nil)
	if err != nil {
		log.Println(err)

		return nil
	}

	for key, value := range Headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)

		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)

		return nil
	}

	var sellers []int

	err = json.Unmarshal(body, &sellers)
	if err != nil {
		log.Println(err)

		return nil
	}

	return sellers
}

func GetOriginalDataSpecialOwnersBySite(key string) []string {
	urlFull := fmt.Sprintf("%s%s", getSpecialOwnersBySiteURL, key)

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, urlFull, nil)
	if err != nil {
		log.Println(err)

		return nil
	}

	for key, value := range Headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)

		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)

		return nil
	}

	var specialOwners []string

	err = json.Unmarshal(body, &specialOwners)
	if err != nil {
		log.Println(err)

		return nil
	}

	return specialOwners
}

func UpdateSpecialOwnersIntoKVS(url string, msg SpecialOwnersMsg) bool {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)

		return false
	}

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)

		return false
	}

	req.Header.Set("Content-Type", "application/json")

	for key, value := range Headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)

		return false
	}
	defer resp.Body.Close()

	// validar si el status code es 200
	if resp.StatusCode == 201 {
		return true
	}

	log.Println(fmt.Sprintf("Error al guardar special owner %s: Status %s error: %s", "extraerKeyDeMSG", resp.Status, resp.Body))

	return false
}
