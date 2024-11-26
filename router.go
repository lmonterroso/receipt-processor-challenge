package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

var ReceiptDB = make(map[string]receipt)

type item struct {
	Description string  `json:"shortDescription"`
	Price       float64 `json:"price,string"`
}

type receipt struct {
	Retailer string `json:"retailer"`
	Date     string `json:"purchaseDate"`
	Time     string `json:"purchaseTime"`
	Total    string `json:"total"`
	Items    []item `json:"items"`
}

type idResp struct {
	ID string `json:"id"`
}

type pointResp struct {
	Points int `json:"points"`
}

func getReceipt(c *gin.Context) {
	id := c.Param("id")
	receipt := ReceiptDB[id]

	points := 0

	//add points for retailer
	for _, char := range receipt.Retailer {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			points++
		}
	}

	//add points for date
	splitDates := strings.Split(receipt.Date, "-")
	if len(splitDates) != 3 {
		c.JSON(http.StatusBadRequest, "Invalid Date format")
		return
	}
	date, err := strconv.Atoi(splitDates[2])
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if date%2 != 0 {
		points += 6
	}

	//add points for time
	if "14:00" < receipt.Time && receipt.Time < "16:00" {
		points += 10
	}

	//add points for total
	parts := strings.Split(receipt.Total, ".")
	if len(parts) != 2 {
		c.JSON(http.StatusBadRequest, "Invalid Price")
		return
	}
	cents, err := strconv.Atoi(parts[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if cents == 0 {
		points += 50
	}
	if cents%25 == 0 {
		points += 25
	}

	//add points for number of items
	points += (len(receipt.Items) / 2) * 5

	for _, item := range receipt.Items {
		trimDesc := strings.TrimSpace(item.Description)
		if len(trimDesc)%3 == 0 {
			roundPoints := item.Price * 0.2

			points += int(math.Ceil(roundPoints))
		}
	}
	c.JSON(http.StatusOK, pointResp{Points: points})
}

func processReceipts(c *gin.Context) {
	var newReceipt receipt
	if err := c.BindJSON(&newReceipt); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	//Marshal newReceipt so we can encode in SHA256
	out, err := json.Marshal(newReceipt)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	h := sha1.New()
	h.Write([]byte(out))
	ReceiptID := hex.EncodeToString(h.Sum(nil))
	fmt.Print(ReceiptID)
	ReceiptDB[ReceiptID] = newReceipt

	c.JSON(http.StatusOK, idResp{ID: ReceiptID})
}
