package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_getPoints(t *testing.T) {
	var server = NewServer()
	example1RequestBody := `{
		"retailer": "Target",
		"purchaseDate": "2022-01-01",
		"purchaseTime": "13:01",
		"items": [
			{
			"shortDescription": "Mountain Dew 12PK",
			"price": "6.49"
			},{
			"shortDescription": "Emils Cheese Pizza",
			"price": "12.25"
			},{
			"shortDescription": "Knorr Creamy Chicken",
			"price": "1.26"
			},{
			"shortDescription": "Doritos Nacho Cheese",
			"price": "3.35"
			},{
			"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
			"price": "12.00"
			}
		],
		"total": "35.35"
	}`

	var example1ReceiptId string

	t.Run("Process receipt example 1", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/receipts/process", strings.NewReader(example1RequestBody))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected status code %v: got %v", http.StatusOK, w.Code)
		}
		var responseJson map[string]string
		err = json.NewDecoder(w.Body).Decode(&responseJson)
		if err != nil {
			t.Errorf("error decoding response body: %v", err)
		}

		var ok bool
		example1ReceiptId, ok = responseJson["id"]
		if !ok {
			t.Errorf("id doesn't exist in response body. response: %v", responseJson)
		}

		r, _ := regexp.Compile(`^\S+$`)
		if !r.MatchString(example1ReceiptId) {
			t.Errorf("receipt id (%v) doesn't match expected format", example1ReceiptId)
		}
	})

	t.Run("Get example 1 receipt points", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/receipts/"+example1ReceiptId+"/points", nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected status code %v: got %v", http.StatusOK, w.Code)
		}

		var responseJson map[string]int
		err = json.NewDecoder(w.Body).Decode(&responseJson)
		if err != nil {
			t.Errorf("error decoding response body (%v) to json: %v", w.Body, err)
		}
		points, ok := responseJson["points"]
		if !ok {
			t.Errorf("points doesn't exist in response body. response: %v", responseJson)
		}

		assert.Equal(t, points, 28)
	})

	t.Run("Get points for invalid receipt id", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/receipts/invalid_receipt_id/points", nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected status code %v: got %v", http.StatusNotFound, w.Code)
		}
	})
}
