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
	prodReaderURL         = "https://production-reader-syncsvc_payment-methods-read-v2.furyapps.io"
	clonReaderURL         = "https://production-reader-comp-clon-readv2.melioffice.com"
	stagingReaderURL      = "https://production-reader-comp-staging-readv2.melioffice.com"
	createStagingURL      = "https://production-synchronizer-stgweb--payment-methods-synchronizer.furyapps.io/v2/payment-methods/golden-gate/refresh"
	createProdAndCloneURL = "https://production-synchronizer-v2--payment-methods-synchronizer.furyapps.io/v2/payment-methods/golden-gate/refresh"
	prodSyncClonReadV2URL = "https://production-synchronizer-clon--payment-methods-read-v2.furyapps.io/pm-core/repository/custom/"
	prodSyncStgReadV2URL  = "https://production-synchronizer-stg--payment-methods-read-v2.furyapps.io/pm-core/repository/custom/"
	prodSyncReadV2URL     = "https://production-synchronizer-v2--payment-methods-read-v2.furyapps.io/pm-core/repository/custom/"
	refreshProdURL        = "https://production-synchronizer-v2--payment-methods-synchronizer.furyapps.io/v2/payment-methods/golden-gate/refresh"
)

var (
	headers = map[string]string{}
)

type InputData struct {
	RefreshProd bool   `json:"refresh_prod"`
	SiteID      string `json:"site_id"`
	SellerIDs   []int  `json:"seller_ids"`
}

type InputSellers struct {
	SellerIDs []int `json:"seller_ids"`
}

type DataResponse struct {
	SiteID          string          `json:"site_id,omitempty"`
	SellerID        int             `json:"seller_id,omitempty"`
	DataAnalysis    *DataAnalysis   `json:"data_analysis,omitempty"`
	ProdData        *SellerAnalysis `json:"prod_data,omitempty"`
	ClonData        *SellerAnalysis `json:"clon_data,omitempty"`
	StagingData     *SellerAnalysis `json:"staging_data,omitempty"`
	OperationDetail string          `json:"operation_detail,omitempty"`
}

type DataAnalysis struct {
	PaymentMethods  PaymentMethodsNode  `json:"payment_methods_node,omitempty"`
	Exclusions      ExclusionsNode      `json:"exclusions_node,omitempty"`
	Groups          GroupsNode          `json:"groups_node,omitempty"`
	AmountAllowed   AmountAllowedNode   `json:"amount_allowed_node,omitempty"`
	OwnPromosByUser OwnPromosByUserNode `json:"own_promos_by_user_node,omitempty"`
}

type PaymentMethodsNode struct {
	Prod    int `json:"prod,omitempty"`
	Clon    int `json:"clon,omitempty"`
	Staging int `json:"stg,omitempty"`
}

type ExclusionsNode struct {
	Prod    int `json:"prod,omitempty"`
	Clon    int `json:"clon,omitempty"`
	Staging int `json:"stg,omitempty"`
}

type GroupsNode struct {
	Prod    int `json:"prod,omitempty"`
	Clon    int `json:"clon,omitempty"`
	Staging int `json:"stg,omitempty"`
}

type AmountAllowedNode struct {
	Prod    int `json:"prod,omitempty"`
	Clon    int `json:"clon,omitempty"`
	Staging int `json:"stg,omitempty"`
}

type OwnPromosByUserNode struct {
	Prod    int `json:"prod,omitempty"`
	Clon    int `json:"clon,omitempty"`
	Staging int `json:"stg,omitempty"`
}

type SellerAnalysis struct {
	ExistsInKVS     bool   `json:"exists_in_kvs,omitempty"`
	SellerID        string `json:"seller_id,omitempty"`
	PaymentMethods  int
	Exclusions      int
	Groups          int
	AmountAllowed   int
	OwnPromosByUser int
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
	CollectorID          string        `json:"collector_id"`
	CustomPaymentMethods []interface{} `json:"custom_payment_methods"`
	Exclusions           []interface{} `json:"exclusions"`
	Groups               interface{}   `json:"groups"`
	AmountAllowed        []interface{} `json:"amount_allowed"`
	OwnPromosByUser      interface{}   `json:"own_promos_by_user"`
}

type OriginalDataKVS struct {
	IsCompressed bool   `json:"IsInStorage"`
	IsInStorage  bool   `json:"IsInStorage"`
	Data         string `json:"Data"`
	LastUpdated  string `json:"LastUpdated"`
}

type LastUpdatedKVS struct {
	SellerID           int    `json:"seller_id"`
	ProdLastUpdated    string `json:"prod_last_updated"`
	ClonLastUpdated    string `json:"clon_last_updated"`
	StagingLastUpdated string `json:"stag_last_updated"`
}

func SetToken(t string) {
	headers = map[string]string{
		"X-Tiger-Token": t,
	}
}

func CompareData(inputData InputData) []DataResponse {
	results := make([]DataResponse, 0)
	refreshProdData(inputData)

	for _, seller := range inputData.SellerIDs {
		sellerResponse := DataResponse{
			SiteID:      inputData.SiteID,
			SellerID:    seller,
			ProdData:    getDataCustomFromURL(prodReaderURL, seller),
			ClonData:    getDataCustomFromURL(clonReaderURL, seller),
			StagingData: getDataCustomFromURL(stagingReaderURL, seller),
		}

		results = append(results, sellerResponse)
	}

	return results
}

func refreshProdData(inputData InputData) {
	if inputData.RefreshProd {
		for _, seller := range inputData.SellerIDs {
			createCustomData(refreshProdURL, inputData.SiteID, seller, "prod")
		}
	}
}

func getDataCustomFromURL(url string, seller int) *SellerAnalysis {
	urlFull := fmt.Sprintf("%s/v2/payment-methods/internal/customizations?collector_id=%d", url, seller)

	client := &http.Client{}

	req, err := http.NewRequest("GET", urlFull, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for key, value := range headers {
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

	var dataCustom DataCustom

	err = json.Unmarshal(body, &dataCustom)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	sellerAnalysis, err := evaluateKVSData(dataCustom)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &sellerAnalysis
}

func evaluateKVSData(dataCustom DataCustom) (SellerAnalysis, error) {
	if len(dataCustom.CustomPaymentMethods) == 0 && len(dataCustom.Exclusions) == 0 &&
		dataCustom.Groups == nil && len(dataCustom.AmountAllowed) == 0 && dataCustom.OwnPromosByUser == nil {

		return SellerAnalysis{}, nil
	}

	sellerAnalysis := SellerAnalysis{
		SellerID:        dataCustom.CollectorID,
		ExistsInKVS:     true,
		PaymentMethods:  len(dataCustom.CustomPaymentMethods),
		Exclusions:      len(dataCustom.Exclusions),
		Groups:          0,
		AmountAllowed:   len(dataCustom.AmountAllowed),
		OwnPromosByUser: 0,
	}

	if dataCustom.Groups != nil {
		sellerAnalysis.Groups = len(dataCustom.Groups.([]interface{}))
	}

	if dataCustom.OwnPromosByUser != nil {
		sellerAnalysis.OwnPromosByUser = len(dataCustom.OwnPromosByUser.([]interface{}))
	}

	return sellerAnalysis, nil
}

func HomologateCustomData(dataComparison []DataResponse) []DataResponse {
	dataResponseList := make([]DataResponse, 0)
	for _, dataSeller := range dataComparison {
		dataResponse := DataResponse{
			SiteID:   dataSeller.SiteID,
			SellerID: dataSeller.SellerID,
		}

		if dataSeller.ProdData.ExistsInKVS == false && dataSeller.ClonData.ExistsInKVS == false && dataSeller.StagingData.ExistsInKVS == false {
			dataResponse.OperationDetail = "No existe en ningún entorno"
		}

		if dataSeller.ProdData.ExistsInKVS == false && dataSeller.ClonData.ExistsInKVS == false && dataSeller.StagingData.ExistsInKVS == true {
			deletedStaging := deleteDataCustomFromURL(prodSyncStgReadV2URL, dataSeller.SellerID)
			dataResponse.OperationDetail = fmt.Sprintf("Eliminar %d en InStaging %s", dataSeller.SellerID, strconv.FormatBool(deletedStaging))
		}

		if dataSeller.ProdData.ExistsInKVS == false && dataSeller.ClonData.ExistsInKVS == true && dataSeller.StagingData.ExistsInKVS == false {
			deletedClon := deleteDataCustomFromURL(prodSyncClonReadV2URL, dataSeller.SellerID)
			dataResponse.OperationDetail = fmt.Sprintf("Eliminar %d en InClon %s", dataSeller.SellerID, strconv.FormatBool(deletedClon))
		}

		if dataSeller.ProdData.ExistsInKVS == false && dataSeller.ClonData.ExistsInKVS == true && dataSeller.StagingData.ExistsInKVS == true {
			deletedClon := deleteDataCustomFromURL(prodSyncClonReadV2URL, dataSeller.SellerID)
			deletedStaging := deleteDataCustomFromURL(prodSyncStgReadV2URL, dataSeller.SellerID)
			dataResponse.OperationDetail = fmt.Sprintf("Eliminar %d en InClon %s y InStaging %s", dataSeller.SellerID, strconv.FormatBool(deletedClon), strconv.FormatBool(deletedStaging))
		}

		if dataSeller.ProdData.ExistsInKVS == true && dataSeller.ClonData.ExistsInKVS == false && dataSeller.StagingData.ExistsInKVS == false {
			createdStaging := createCustomData(createStagingURL, dataSeller.SiteID, dataSeller.SellerID, "staging")
			createdProd := createCustomData(createProdAndCloneURL, dataSeller.SiteID, dataSeller.SellerID, "prodAndClone")
			dataResponse.OperationDetail = fmt.Sprintf("Se creó Prod %s and InStaging %s %d", strconv.FormatBool(createdProd), strconv.FormatBool(createdStaging), dataSeller.SellerID)
		}

		if dataSeller.ProdData.ExistsInKVS == true && dataSeller.ClonData.ExistsInKVS == false && dataSeller.StagingData.ExistsInKVS == true {
			dataResponse.OperationDetail = fmt.Sprintf("Crear en InClon %d", dataSeller.SellerID)
		}

		if dataSeller.ProdData.ExistsInKVS == true && dataSeller.ClonData.ExistsInKVS == true && dataSeller.StagingData.ExistsInKVS == false {
			created := createCustomData(createStagingURL, dataSeller.SiteID, dataSeller.SellerID, "staging")
			dataResponse.OperationDetail = fmt.Sprintf("Se creó InStaging %d %s", dataSeller.SellerID, strconv.FormatBool(created))
		}

		if dataSeller.ProdData.ExistsInKVS == true && dataSeller.ClonData.ExistsInKVS == true && dataSeller.StagingData.ExistsInKVS == true {
			dataResponse.DataAnalysis = &DataAnalysis{
				PaymentMethods: PaymentMethodsNode{
					Prod:    dataSeller.ProdData.PaymentMethods,
					Clon:    dataSeller.ClonData.PaymentMethods,
					Staging: dataSeller.StagingData.PaymentMethods,
				},
				Exclusions: ExclusionsNode{
					Prod:    dataSeller.ProdData.Exclusions,
					Clon:    dataSeller.ClonData.Exclusions,
					Staging: dataSeller.StagingData.Exclusions,
				},
				Groups: GroupsNode{
					Prod:    dataSeller.ProdData.Groups,
					Clon:    dataSeller.ClonData.Groups,
					Staging: dataSeller.StagingData.Groups,
				},
				AmountAllowed: AmountAllowedNode{
					Prod:    dataSeller.ProdData.AmountAllowed,
					Clon:    dataSeller.ClonData.AmountAllowed,
					Staging: dataSeller.StagingData.AmountAllowed,
				},
				OwnPromosByUser: OwnPromosByUserNode{
					Prod:    dataSeller.ProdData.OwnPromosByUser,
					Clon:    dataSeller.ClonData.OwnPromosByUser,
					Staging: dataSeller.StagingData.OwnPromosByUser,
				},
			}
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

func ValidateLastUpdateIntoKVS(inputSellers InputSellers) []LastUpdatedKVS {
	lastUpdatedKVS := make([]LastUpdatedKVS, 0)
	for _, seller := range inputSellers.SellerIDs {
		sellerData := LastUpdatedKVS{
			SellerID:           seller,
			ProdLastUpdated:    getOriginalDataCustomFromURL(prodSyncReadV2URL, seller).LastUpdated,
			ClonLastUpdated:    getOriginalDataCustomFromURL(prodSyncClonReadV2URL, seller).LastUpdated,
			StagingLastUpdated: getOriginalDataCustomFromURL(prodSyncStgReadV2URL, seller).LastUpdated,
		}
		lastUpdatedKVS = append(lastUpdatedKVS, sellerData)
	}

	return lastUpdatedKVS
}

func getOriginalDataCustomFromURL(url string, seller int) *OriginalDataKVS {
	urlFull := fmt.Sprintf("%s%d", url, seller)

	client := &http.Client{}

	req, err := http.NewRequest("GET", urlFull, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for key, value := range headers {
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

	var originalDataKVS OriginalDataKVS

	err = json.Unmarshal(body, &originalDataKVS)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &originalDataKVS
}
