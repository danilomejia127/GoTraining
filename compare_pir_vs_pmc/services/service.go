package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// Scope para eliminación production-synchronizer-stg--payment-methods-read-v2.furyapps.io
const (
	createStagingURL      = "https://production-synchronizer-stgweb--payment-methods-synchronizer.furyapps.io/v2/payment-methods/golden-gate/refresh"
	createProdAndCloneURL = "https://production-synchronizer-v2--payment-methods-synchronizer.furyapps.io/v2/payment-methods/golden-gate/refresh"
	deleteClonURL         = "https://production-synchronizer-clon--payment-methods-read-v2.furyapps.io/pm-core/repository/custom/"
	deleteStagingURL      = "https://production-synchronizer-stg--payment-methods-read-v2.furyapps.io/pm-core/repository/custom/"
)

var (
	headers = map[string]string{}
)

type InputData struct {
	SiteID    string `json:"site_id"`
	SellerIDs []int  `json:"seller_ids"`
}

type DataComparition struct {
	SiteID    string
	SellerID  int
	InProd    bool
	InClon    bool
	InStaging bool
}

type DataResponse struct {
	SiteID          string
	SellerID        int
	OperationDetail string
}

type RefreshMessage struct {
	Msg struct {
		ID struct {
			SiteID   string `json:"site_id"`
			SellerID string `json:"seller_id"`
		} `json:"id"`
	} `json:"msg"`
}

type DataCustom struct {
	CollectorID          string `json:"collector_id"`
	CustomPaymentMethods []any  `json:"custom_payment_methods"`
	Exclusions           []any  `json:"exclusions"`
	Groups               any    `json:"groups"`
	AmountAllowed        []any  `json:"amount_allowed"`
	OwnPromosByUser      any    `json:"own_promos_by_user"`
}

func SetToken(t string) {
	headers = map[string]string{
		"X-Tiger-Token": t,
	}
}

func CompareData(inputData InputData) []DataComparition {
	prodURL := "https://production-reader-syncsvc_payment-methods-read-v2.furyapps.io"
	clonURL := "https://production-reader-comp-clon-readv2.melioffice.com"
	stagingURL := "https://production-reader-comp-staging-readv2.melioffice.com"

	results := make([]DataComparition, 0)

	for _, seller := range inputData.SellerIDs {
		sellerResponse := DataComparition{
			SiteID:    inputData.SiteID,
			SellerID:  seller,
			InProd:    getDataCustomFromURL(prodURL, seller),
			InClon:    getDataCustomFromURL(clonURL, seller),
			InStaging: getDataCustomFromURL(stagingURL, seller),
		}

		results = append(results, sellerResponse)
	}

	return results
}

func getDataCustomFromURL(url string, seller int) bool {
	urlFull := fmt.Sprintf("%s/v2/payment-methods/internal/customizations?collector_id=%d", url, seller)

	client := &http.Client{}

	req, err := http.NewRequest("GET", urlFull, nil)
	if err != nil {
		fmt.Println(err)
		return false
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}

	var dataCustom DataCustom

	err = json.Unmarshal(body, &dataCustom)
	if err != nil {
		fmt.Println(err)
		return false
	}

	result, err := isDataInKVS(dataCustom)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return result
}

func isDataInKVS(dataCustom DataCustom) (bool, error) {
	if len(dataCustom.CustomPaymentMethods) == 0 && len(dataCustom.Exclusions) == 0 &&
		dataCustom.Groups == nil && len(dataCustom.AmountAllowed) == 0 && dataCustom.OwnPromosByUser == nil {

		return false, nil
	}

	return true, nil
}

func HomologateCustomData(dataComparition []DataComparition) []DataResponse {
	dataResponseList := make([]DataResponse, 0)
	for _, dataSeller := range dataComparition {
		dataResponse := DataResponse{
			SiteID:          dataSeller.SiteID,
			SellerID:        dataSeller.SellerID,
			OperationDetail: "Existe en los tres entornos",
		}

		if dataSeller.InProd == false && dataSeller.InClon == false && dataSeller.InStaging == false {
			dataResponse.OperationDetail = "No existe en ningún entorno"
		}

		if dataSeller.InProd == false && dataSeller.InClon == false && dataSeller.InStaging == true {
			deletedStaging := deleteDataCustomFromURL(deleteStagingURL, dataSeller.SellerID)
			dataResponse.OperationDetail = fmt.Sprintf("Eliminar %d en InStaging %s", dataSeller.SellerID, strconv.FormatBool(deletedStaging))
		}

		if dataSeller.InProd == false && dataSeller.InClon == true && dataSeller.InStaging == false {
			deletedClon := deleteDataCustomFromURL(deleteClonURL, dataSeller.SellerID)
			dataResponse.OperationDetail = fmt.Sprintf("Eliminar %d en InClon %s", dataSeller.SellerID, strconv.FormatBool(deletedClon))
		}

		if dataSeller.InProd == false && dataSeller.InClon == true && dataSeller.InStaging == true {
			deletedClon := deleteDataCustomFromURL(deleteClonURL, dataSeller.SellerID)
			deletedStaging := deleteDataCustomFromURL(deleteStagingURL, dataSeller.SellerID)
			dataResponse.OperationDetail = fmt.Sprintf("Eliminar %d en InClon %s y InStaging %s", dataSeller.SellerID, strconv.FormatBool(deletedClon), strconv.FormatBool(deletedStaging))
		}

		if dataSeller.InProd == true && dataSeller.InClon == false && dataSeller.InStaging == false {
			createdStaging := createCustomData(createStagingURL, dataSeller.SiteID, dataSeller.SellerID, "staging")
			createdProd := createCustomData(createProdAndCloneURL, dataSeller.SiteID, dataSeller.SellerID, "prodAndClone")
			dataResponse.OperationDetail = fmt.Sprintf("Se creó Prod %s and InStaging %s %d", strconv.FormatBool(createdProd), strconv.FormatBool(createdStaging), dataSeller.SellerID)
		}

		if dataSeller.InProd == true && dataSeller.InClon == false && dataSeller.InStaging == true {
			dataResponse.OperationDetail = fmt.Sprintf("Crear en InClon %d", dataSeller.SellerID)
		}

		if dataSeller.InProd == true && dataSeller.InClon == true && dataSeller.InStaging == false {
			created := createCustomData(createStagingURL, dataSeller.SiteID, dataSeller.SellerID, "staging")
			dataResponse.OperationDetail = fmt.Sprintf("Se creó InStaging %d %s", dataSeller.SellerID, strconv.FormatBool(created))
		}

		dataResponseList = append(dataResponseList, dataResponse)
	}

	return dataResponseList
}

func createCustomData(refreshURL string, siteID string, sellerID int, scope string) bool {
	msg := RefreshMessage{}
	msg.Msg.ID.SiteID = siteID
	msg.Msg.ID.SellerID = strconv.Itoa(sellerID)

	jsonData, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return false
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", refreshURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
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

func deleteDataCustomFromURL(url string, seller int) bool {
	urlFull := url + fmt.Sprintf("%d", seller)

	client := &http.Client{}

	req, err := http.NewRequest("DELETE", urlFull, nil)
	if err != nil {
		fmt.Println(err)
		return false
	}

	for key, value := range headers {
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
