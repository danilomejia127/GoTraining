package services

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/mercadolibre/GoTraining/compare_pir_vs_pmc/apicalls"
	"github.com/mercadolibre/GoTraining/compare_pir_vs_pmc/dtos"
)

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
	ProdData        *dtos.Collector `json:"prod_data,omitempty"`
	ClonData        *dtos.Collector `json:"clon_data,omitempty"`
	StagingData     *dtos.Collector `json:"staging_data,omitempty"`
	OperationDetail string          `json:"operation_detail,omitempty"`
}

type Response struct {
	DataResponse    []DataResponse `json:"data,omitempty"`
	SellersWithData []int          `json:"seller_ids"`
}

type DataAnalysis struct {
	PaymentMethods  PaymentMethodsNode  `json:"payment_methods_node,omitempty"`
	Exclusions      ExclusionsNode      `json:"exclusions_node,omitempty"`
	Groups          GroupsNode          `json:"groups_node,omitempty"`
	AmountAllowed   AmountAllowedNode   `json:"amount_allowed_node,omitempty"`
	OwnPromosByUser OwnPromosByUserNode `json:"own_promos_by_user_node,omitempty"`
}

type PaymentMethodsNode struct {
	Prod               int                 `json:"prod,omitempty"`
	Clon               int                 `json:"clon,omitempty"`
	Staging            int                 `json:"stg,omitempty"`
	PaymentMethodsDiff *PaymentMethodsDiff `json:"payment_methods_diff,omitempty"`
}

type PaymentMethodsDiff struct {
	Message       string `json:"message,omitempty"`
	IdsOnlyInProd []int  `json:"ids_only_in_prod,omitempty"`
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

type LastUpdatedKVS struct {
	SellerID           int    `json:"seller_id"`
	ProdLastUpdated    string `json:"prod_last_updated"`
	ClonLastUpdated    string `json:"clon_last_updated"`
	StagingLastUpdated string `json:"stag_last_updated"`
}

func CompareData(inputData InputData) Response {
	results := make([]DataResponse, 0)
	refreshProdData(inputData)

	for _, seller := range inputData.SellerIDs {
		sellerResponse := DataResponse{
			SiteID:      inputData.SiteID,
			SellerID:    seller,
			ProdData:    apicalls.GetDataCustomFromURL(prodReaderURL, seller),
			ClonData:    apicalls.GetDataCustomFromURL(clonReaderURL, seller),
			StagingData: apicalls.GetDataCustomFromURL(stagingReaderURL, seller),
		}

		results = append(results, sellerResponse)
	}

	return Response{
		DataResponse: results,
	}
}

func refreshProdData(inputData InputData) {
	if inputData.RefreshProd {
		var wg sync.WaitGroup
		semaphore := make(chan struct{}, 10)

		for _, seller := range inputData.SellerIDs {
			wg.Add(1)
			go func(sellerID int) { // Pasar seller como parámetro
				defer wg.Done()
				semaphore <- struct{}{}
				log.Println("Refreshing seller " + strconv.Itoa(sellerID))
				apicalls.CreateCustomData(refreshProdURL, inputData.SiteID, sellerID, "prod")
				apicalls.CreateCustomData(createStagingURL, inputData.SiteID, sellerID, "stag")
				<-semaphore
			}(seller) // Pasar el valor actual de seller como argumento
		}

		wg.Wait()
	}
}

func existsCustomSeller(dataCustom *dtos.Collector) bool {
	if dataCustom != nil && len(dataCustom.PaymentMethods) == 0 && len(dataCustom.Exclusions) == 0 &&
		dataCustom.Groups == nil && len(dataCustom.AmountAllowed) == 0 && dataCustom.OwnPromosByUser == nil {

		return false
	}

	return true
}

func HomologateCustomData(dataComparison Response) Response {
	dataResponseList := make([]DataResponse, 0)
	sellersWithData := make([]int, 0)

	for _, dataSeller := range dataComparison.DataResponse {
		if dataSeller.ProdData == nil && dataSeller.ClonData == nil && dataSeller.StagingData == nil {
			continue
		}

		existsInPrd := existsCustomSeller(dataSeller.ProdData)
		existsInClon := existsCustomSeller(dataSeller.ClonData)
		existsInStg := existsCustomSeller(dataSeller.StagingData)

		dataResponse := DataResponse{
			SiteID:   dataSeller.SiteID,
			SellerID: dataSeller.SellerID,
		}

		if existsInPrd == false && existsInClon == false && existsInStg == false {
			dataResponse.OperationDetail = "No existe en ningún entorno"
		}

		if existsInPrd == false && existsInClon == false && existsInStg == true {
			deletedStaging := apicalls.DeleteDataCustomFromURL(prodSyncStgReadV2URL, dataSeller.SellerID)
			dataResponse.OperationDetail = fmt.Sprintf("Eliminar %d en InStaging %s", dataSeller.SellerID, strconv.FormatBool(deletedStaging))
		}

		if existsInPrd == false && existsInClon == true && existsInStg == false {
			deletedClon := apicalls.DeleteDataCustomFromURL(prodSyncClonReadV2URL, dataSeller.SellerID)
			dataResponse.OperationDetail = fmt.Sprintf("Eliminar %d en InClon %s", dataSeller.SellerID, strconv.FormatBool(deletedClon))
		}

		if existsInPrd == false && existsInClon == true && existsInStg == true {
			deletedClon := apicalls.DeleteDataCustomFromURL(prodSyncClonReadV2URL, dataSeller.SellerID)
			deletedStaging := apicalls.DeleteDataCustomFromURL(prodSyncStgReadV2URL, dataSeller.SellerID)
			dataResponse.OperationDetail = fmt.Sprintf("Eliminar %d en InClon %s y InStaging %s", dataSeller.SellerID, strconv.FormatBool(deletedClon), strconv.FormatBool(deletedStaging))
		}

		if existsInPrd == true && existsInClon == false && existsInStg == false {
			createdStaging := apicalls.CreateCustomData(createStagingURL, dataSeller.SiteID, dataSeller.SellerID, "staging")
			createdProd := apicalls.CreateCustomData(createProdAndCloneURL, dataSeller.SiteID, dataSeller.SellerID, "prodAndClone")
			dataResponse.OperationDetail = fmt.Sprintf("Se creó Prod %s and InStaging %s %d", strconv.FormatBool(createdProd), strconv.FormatBool(createdStaging), dataSeller.SellerID)
		}

		if existsInPrd == true && existsInClon == false && existsInStg == true {
			dataResponse.OperationDetail = fmt.Sprintf("Crear en InClon %d", dataSeller.SellerID)
		}

		if existsInPrd == true && existsInClon == true && existsInStg == false {
			created := apicalls.CreateCustomData(createStagingURL, dataSeller.SiteID, dataSeller.SellerID, "staging")
			dataResponse.OperationDetail = fmt.Sprintf("Se creó InStaging %d %s", dataSeller.SellerID, strconv.FormatBool(created))
		}

		if existsInPrd == true && existsInClon == true && existsInStg == true {
			lenProd, lenClon, lenStg, paymentMethodsDiff := validatePaymentMethodsNode(dataSeller)

			dataResponse.DataAnalysis = &DataAnalysis{
				PaymentMethods: PaymentMethodsNode{
					Prod:               lenProd,
					Clon:               lenClon,
					Staging:            lenStg,
					PaymentMethodsDiff: paymentMethodsDiff,
				},
				Exclusions: ExclusionsNode{
					Prod:    len(dataSeller.ProdData.Exclusions),
					Clon:    len(dataSeller.ClonData.Exclusions),
					Staging: len(dataSeller.StagingData.Exclusions),
				},
				Groups: GroupsNode{
					Prod:    getLengthOfNode(dataSeller.ProdData.Groups),
					Clon:    getLengthOfNode(dataSeller.ClonData.Groups),
					Staging: getLengthOfNode(dataSeller.StagingData.Groups),
				},
				AmountAllowed: AmountAllowedNode{
					Prod:    len(dataSeller.ProdData.AmountAllowed),
					Clon:    len(dataSeller.ClonData.AmountAllowed),
					Staging: len(dataSeller.StagingData.AmountAllowed),
				},
				OwnPromosByUser: OwnPromosByUserNode{
					Prod:    getLengthOfNode(dataSeller.ProdData.OwnPromosByUser),
					Clon:    getLengthOfNode(dataSeller.ClonData.OwnPromosByUser),
					Staging: getLengthOfNode(dataSeller.StagingData.OwnPromosByUser),
				},
			}

			sellersWithData = append(sellersWithData, dataSeller.SellerID)
		}

		dataResponseList = append(dataResponseList, dataResponse)
	}

	return Response{
		DataResponse:    dataResponseList,
		SellersWithData: sellersWithData,
	}
}

func getLengthOfNode(node interface{}) int {
	if node == nil {
		return 0
	}

	return len(node.([]interface{}))
}

func ValidateLastUpdateIntoKVS(inputSellers InputSellers) []LastUpdatedKVS {
	lastUpdatedKVS := make([]LastUpdatedKVS, 0)

	for _, seller := range inputSellers.SellerIDs {
		sellerData := LastUpdatedKVS{
			SellerID:           seller,
			ProdLastUpdated:    apicalls.GetOriginalDataCustomFromURL(prodSyncReadV2URL, seller).LastUpdated,
			ClonLastUpdated:    apicalls.GetOriginalDataCustomFromURL(prodSyncClonReadV2URL, seller).LastUpdated,
			StagingLastUpdated: apicalls.GetOriginalDataCustomFromURL(prodSyncStgReadV2URL, seller).LastUpdated,
		}
		lastUpdatedKVS = append(lastUpdatedKVS, sellerData)
	}

	return lastUpdatedKVS
}

func validatePaymentMethodsNode(dataSeller DataResponse) (int, int, int, *PaymentMethodsDiff) {
	lenProd := len(dataSeller.ProdData.PaymentMethods)
	lenClon := len(dataSeller.ClonData.PaymentMethods)
	lenStg := len(dataSeller.StagingData.PaymentMethods)

	if lenClon != lenStg {
		return lenProd, lenClon, lenStg, &PaymentMethodsDiff{
			IdsOnlyInProd: findMissingPaymentMethods(dataSeller.ProdData.PaymentMethods, dataSeller.StagingData.PaymentMethods),
		}
	}

	return lenProd, lenClon, lenStg, nil
}

func findMissingPaymentMethods(dataProd, staging []dtos.PaymentMethod) []int {
	missingPaymentMethods := make([]int, 0)

	stagingMap := make(map[string]bool)
	for _, paymentMethod := range staging {
		stagingMap[paymentMethod.Misc.PmIssuerRelation.ID] = true
	}

	for _, paymentMethod := range dataProd {
		// Verificar si el medio de pago no existe en staging
		if !stagingMap[paymentMethod.Misc.PmIssuerRelation.ID] {
			pmID, _ := strconv.Atoi(paymentMethod.Misc.PmIssuerRelation.ID)
			missingPaymentMethods = append(missingPaymentMethods, pmID)
		}
	}

	return missingPaymentMethods
}
