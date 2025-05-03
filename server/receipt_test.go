package server

import (
	"context"
	"testing"
)

type receiptWithPoints struct {
	receipt Receipt
	points  int64
}

var test_receipts map[string]receiptWithPoints = map[string]receiptWithPoints{
	"dollar-general": {Receipt{"Dollar General", "2022-01-01", "13:01", []Item{{"abc", "10.00"}, {"abcd", "20.39"}}, "30.39"}, 26},
	// dollar-general points: 13 + 5 + 2 + 6 = 26
	"walmart": {Receipt{"Walmart1", "2022-01-02", "14:01", []Item{{"   abc   ", "6.49"}, {"abcd  ", "1.26"}}, "7.75"}, 50},
	// walmart points: 8 + 25 + 5 + 2 + 10 = 50
	"readme-example1": {Receipt{
		"Target",
		"2022-01-01",
		"13:01",
		[]Item{
			{"Mountain Dew 12PK", "6.49"},
			{"Emils Cheese Pizza", "12.25"},
			{"Knorr Creamy Chicken", "1.26"},
			{"Doritos Nacho Cheese", "3.35"},
			{"Klarbrunn 12-PK 12 FL OZ", "12.00"},
		},
		"35.35",
	}, 28},
	"readme-example2": {Receipt{
		"M&M Corner Market",
		"2022-03-20",
		"14:33",
		[]Item{
			{"Gatorade", "2.25"},
			{"Gatorade", "2.25"},
			{"Gatorade", "2.25"},
			{"Gatorade", "2.25"},
		},
		"9.00",
	}, 109},
}

func TestItem_PriceValue(t *testing.T) {
	type priceValueTest struct {
		item Item
		want float64
	}

	tests := []priceValueTest{
		{Item{"Good 1", "10.00"}, 10},
		{Item{"Good 2", "20.39"}, 20.39},
		{Item{"Good 3", "0.00"}, 0},
		{Item{"Bad", "abc"}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.item.ShortDescription, func(t *testing.T) {
			if got := tt.item.PriceValue(context.Background()); got != tt.want {
				t.Errorf("Item.PriceValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReceipt_total(t *testing.T) {
	type totalTest struct {
		receipt Receipt
		want    float64
	}

	tests := []totalTest{
		{Receipt{Total: "30.39"}, 30.39},
		{Receipt{Total: "abc"}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.receipt.Total, func(t *testing.T) {
			if got := tt.receipt.total(context.Background()); got != tt.want {
				t.Errorf("Receipt.total() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReceipt_Points(t *testing.T) {
	for name, receipt := range test_receipts {
		t.Run(name, func(t *testing.T) {
			if got := receipt.receipt.Points(context.Background()); got != receipt.points {
				t.Errorf("Receipt.Points() = %v, want %v", got, receipt.points)
			}
		})
	}
}
