package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

const (
	runtime = 10 * time.Minute    // Duration of each phase 	// Number of concurrent goroutines
	token   = "Bearer here token" // Token to add in the header
)

var (
	siteIds  = []string{"MLA"}
	idValues = []string{
		"account_money,pagofacil,rapipago,debin_transfer", // With all IDs
	}
)

// Job structure to hold request info
type Job struct {
	siteID    string
	id        string
	addHeader bool
}

type StandardSources struct {
	PaymentMethods []interface{} `json:"results"`
	Paging         struct {      // Deprecated?
		Total  int `json:"total"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	} `json:"paging"`
}

func sendRequest(job Job) {
	url := fmt.Sprintf("https://testing-reader-data--payment-methods-read-v2.furyapps.io/v1/payment_methods/search?site_id=%s&marketplace=MELI&status=active&product_id=be060sp2le1g01lpjqj0&id=%s&caller.id=221665156&attributes=id,name,status,payment_method_id,secure_thumbnail,thumbnail,payment_type_id,issuer,payer_costs,deferred_capture,min_accreditation_days,max_accreditation_days,accreditation_time,assets", job.siteID, job.id)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating the request: %s\n", err)
		return
	}

	// Add the x-tiger-token header
	req.Header.Add("x-tiger-token", token)

	// Add the X-Core-Flow-Type header if required
	if job.addHeader {
		fmt.Println("Adding header X-Core-Flow-Type")
		req.Header.Add("X-Core-Flow-Type", "buyingflow-pm-off")
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error in the request: %s\n", err)

		return
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Non-200 response: %d, body: %s\n", resp.StatusCode, body)

		return
	}

	paymentMethodResponse := &StandardSources{}

	err = json.NewDecoder(resp.Body).Decode(paymentMethodResponse)
	if err != nil {
		fmt.Printf("Error decoding the response: %s\n", err)

		return
	}

	fmt.Printf("Response for site_id=%s, id=%s, header X-Core-Flow-Type: %t Status: %s total: %d \n", job.siteID, job.id, job.addHeader, resp.Status, paymentMethodResponse.Paging.Total)
}

func runRequests(duration time.Duration, id string, requestCounter *int, withHeader bool) {
	startTime := time.Now()

	// Generate requests during the specified time
	for time.Since(startTime) < duration {
		siteID := siteIds[rand.Intn(len(siteIds))]
		sendRequest(Job{siteID, id, withHeader})
		time.Sleep(20 * time.Millisecond) // Sleep to avoid overwhelming the server

		(*requestCounter)++ // Count the request
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Counters for each type of request
	var requestsWithoutID int
	var requestsWithID int
	var requestsWithoutIDWithHeader int

	// ** Phase 1: Request without ID (no headers) for 5 minutes **
	//runRequests(runtime, "", &requestsWithoutID, false)
	//time.Sleep(2 * time.Second)

	// ** Phase 2: Request with IDs (no headers) for 5 minutes **
	//runRequests(runtime, idValues[0], &requestsWithID, false)
	//time.Sleep(2 * time.Second)

	// ** Phase 3: Request without ID (with header) for 5 minutes **
	runRequests(runtime, "", &requestsWithoutIDWithHeader, true)

	fmt.Printf("Total requests with ID: %d\n", requestsWithID)
	fmt.Printf("Total requests without ID (no header): %d\n", requestsWithoutID)
	fmt.Printf("Total requests without ID (with header): %d\n", requestsWithoutIDWithHeader)

	fmt.Println("All requests have been processed.")
}
