package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Scope para eliminación production-synchronizer-stg--payment-methods-read-v2.furyapps.io

type DataComparition struct {
	SellerId  int
	InProd    bool
	InClon    bool
	InStaging bool
}

type DataCustom struct {
	CollectorID          string `json:"collector_id"`
	CustomPaymentMethods []any  `json:"custom_payment_methods"`
	Exclusions           []any  `json:"exclusions"`
	Groups               any    `json:"groups"`
	AmountAllowed        []any  `json:"amount_allowed"`
	OwnPromosByUser      any    `json:"own_promos_by_user"`
}

func CompareData(sellers []int) []DataComparition {
	token := "Bearer "
	prodURL := "https://production-reader-syncsvc_payment-methods-read-v2.furyapps.io"
	clonURL := "https://production-reader-comp-clon-readv2.melioffice.com"
	stagingURL := "https://production-reader-comp-staging-readv2.melioffice.com"

	results := make([]DataComparition, 0)

	for _, seller := range sellers {
		sellerResponse := DataComparition{
			SellerId:  seller,
			InProd:    getDataCustomFromURL(prodURL, seller, token),
			InClon:    getDataCustomFromURL(clonURL, seller, token),
			InStaging: getDataCustomFromURL(stagingURL, seller, token),
		}

		results = append(results, sellerResponse)
	}

	return results
}

func getDataCustomFromURL(url string, seller int, token string) bool {
	urlFull := fmt.Sprintf("%s/v2/payment-methods/internal/customizations?collector_id=%d", url, seller)

	headers := map[string]string{
		"X-Tiger-Token": token,
	}

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

func HomologateCustomData(dataComparition []DataComparition) {
	for _, dataSeller := range dataComparition {
		if dataSeller.InProd == false && dataSeller.InClon == false && dataSeller.InStaging == true {
			fmt.Println(fmt.Sprintf("Eliminar en InStaging %d", dataSeller.SellerId))
		}

		if dataSeller.InProd == false && dataSeller.InClon == true && dataSeller.InStaging == false {
			fmt.Println(fmt.Sprintf("Eliminar en InClon %d", dataSeller.SellerId))
		}

		if dataSeller.InProd == false && dataSeller.InClon == true && dataSeller.InStaging == true {
			fmt.Println(fmt.Sprintf("Eliminar en InClon y InStaging %d", dataSeller.SellerId))
		}

		if dataSeller.InProd == true && dataSeller.InClon == false && dataSeller.InStaging == false {
			fmt.Println(fmt.Sprintf("Crear en InClon y InStaging %d", dataSeller.SellerId))
		}

		if dataSeller.InProd == true && dataSeller.InClon == false && dataSeller.InStaging == true {
			fmt.Println(fmt.Sprintf("Crear en InClon %d", dataSeller.SellerId))
		}

		if dataSeller.InProd == true && dataSeller.InClon == true && dataSeller.InStaging == false {
			fmt.Println(fmt.Sprintf("Crear en InStaging %d", dataSeller.SellerId))
		}
	}
}
