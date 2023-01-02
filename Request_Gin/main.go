package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
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

	for _, seller := range siteSellerID.Sellers {
		fmt.Println(seller)
	}
	return
}
