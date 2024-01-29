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
	results := make([]DataComparition, 0)

	for _, seller := range sellers {
		sellerResponse := DataComparition{
			SellerId: seller,
			InProd:   getProdCustom(seller),
			InClon:   getClonCustom(seller),
		}

		results = append(results, sellerResponse)
	}

	return results
}

func getProdCustom(seller int) bool {
	prodURL := fmt.Sprintf("https://production-reader-syncsvc_payment-methods-read-v2.furyapps.io/v2/payment-methods/internal/customizations?collector_id=%d", seller)

	headers := map[string]string{
		"X-Tiger-Token": "Bearer ",
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", prodURL, nil)
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

	result, err := isEmptyDataCustom(dataCustom)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return result
}

func getClonCustom(seller int) bool {

	return false
}

func isEmptyDataCustom(dataCustom DataCustom) (bool, error) {
	if len(dataCustom.CustomPaymentMethods) == 0 && len(dataCustom.Exclusions) == 0 &&
		dataCustom.Groups == nil && len(dataCustom.AmountAllowed) == 0 && dataCustom.OwnPromosByUser == nil {

		return true, nil
	}

	return false, nil
}
