package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"time"
)

var (
	apiUrl   = "https://my-scope"
	resource = "/v2/payment-methods/golden-gate/refresh"
)

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSONP(200, gin.H{
			"message": "pong",
		})
	})
	router.POST("create_messages", createMessages)
	router.Run("localhost:8080")
}

type RefreshGoldenGate struct {
	Msg struct {
		ID struct {
			SellerID string `json:"seller_id"`
			SiteID   string `json:"site_id"`
		} `json:"id"`
	} `json:"msg"`
}

type SiteSellerID struct {
	Site    string
	Sellers []string
}

func createMessages(c *gin.Context) {

	var siteSellerID SiteSellerID

	if err := c.BindJSON(&siteSellerID); err != nil {
		c.JSONP(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	client := &http.Client{}
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()
	var buf bytes.Buffer
	siteID := siteSellerID.Site

	for _, seller := range siteSellerID.Sellers {
		rgg := RefreshGoldenGate{}
		rgg.Msg.ID.SiteID = siteID
		rgg.Msg.ID.SellerID = seller
		err := json.NewEncoder(&buf).Encode(rgg)
		if err != nil {
			c.JSONP(500, gin.H{
				"message": err.Error(),
			})
			return
		}

		req, err := http.NewRequest(http.MethodPost, urlStr, &buf)
		req.Header.Add("X-Tiger-Token", "Token_Here")
		if err != nil {
			c.JSONP(500, gin.H{
				"message": err.Error(),
			})
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			c.JSONP(500, gin.H{
				"message": err.Error(),
			})
			return
		}

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSONP(500, gin.H{
				"message": err.Error(),
			})
			return
		}

		fmt.Println(resp.Status)
		fmt.Println(seller)
		fmt.Println(string(bytes))
		time.Sleep(100 * time.Millisecond)
	}
	return
}
