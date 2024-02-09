package apicalls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
