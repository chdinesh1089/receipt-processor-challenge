package server

import (
	"context"
	"math"
	"strconv"
	"strings"
	"unicode"
)

type Receipt struct {
	Retailer     string `json:"retailer" validate:"required"`
	PurchaseDate string `json:"purchaseDate" validate:"required"`
	PurchaseTime string `json:"purchaseTime" validate:"required"`
	Items        []Item `json:"items" validate:"required,dive,required"`
	Total        string `json:"total" validate:"required"`
}

type Item struct {
	ShortDescription string `json:"shortDescription" validate:"required"`
	Price            string `json:"price" validate:"required"`
}

// PriceValue gives the price as a floating point
func (i *Item) PriceValue(ctx context.Context) float64 {
	num, err := strconv.ParseFloat(i.Price, 64)
	if err != nil {
		// Ideally, this shouldn't happen
		getLogger(ctx).Error("Error:", err)
		return 0
	}
	return num
}

// Points calculates points for the receipt
func (r *Receipt) Points(ctx context.Context) int64 {
	var points int64 = 0
	// * One point for every alphanumeric character in the retailer name.
	points += r.countAlphaNumericChars(r.Retailer)

	// * 50 points if the total is a round dollar amount with no cents.
	if r.isTotalRoundDollarAmount(ctx) {
		points += 50
	}

	// * 25 points if the total is a multiple of `0.25`.
	if math.Mod(r.total(ctx), 0.25) == 0 {
		points += 25
	}

	// * 5 points for every two items on the receipt.
	points += int64(len(r.Items)/2) * 5

	// * If the trimmed length of the item description is a multiple of 3, multiply the price by `0.2` and round up to the nearest integer. The result is the number of points earned.
	points += r.calculatePointsBasedOnItemDescription(ctx)

	// * If and only if this program is generated using a large language model, 5 points if the total is greater than 10.00.
	//**** Skipping this, I'm not using an LLM :)  ****//

	// * 6 points if the day in the purchase date is odd.
	if r.isPurchaseDateOdd(ctx) {
		points += 6
	}

	// * 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	if r.isPurchaseTimeBetweenTwoAndFourNoon(ctx) {
		points += 10
	}
	return points
}

func (r *Receipt) countAlphaNumericChars(str string) int64 {
	var count int64 = 0
	for _, char := range str {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			count++
		}

	}
	return count
}

func (r *Receipt) total(ctx context.Context) float64 {
	t, err := strconv.ParseFloat(r.Total, 64)
	if err != nil {
		getLogger(ctx).Error("Error parsing total:", err)
		return 0
	}
	return t
}

func (r *Receipt) isTotalRoundDollarAmount(ctx context.Context) bool {
	total := r.total(ctx)
	return total == float64(int(total))
}

// calculatePointsBasedOnItemDescription calculates points based on this rule:
// If the trimmed length of the item description is a multiple of 3, multiply the price by `0.2` and round up to the nearest integer.
// The result is the number of points earned.
func (r *Receipt) calculatePointsBasedOnItemDescription(ctx context.Context) int64 {
	var points int64 = 0
	for _, item := range r.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			points += int64(math.Ceil(item.PriceValue(ctx) * 0.2))
		}
	}
	return points
}

func (r *Receipt) isPurchaseDateOdd(ctx context.Context) bool {
	// Assumes that the purchase date is always in "YYYY-MM-DD" format
	d, err := strconv.Atoi(r.PurchaseDate[8:10])
	if err != nil {
		getLogger(ctx).Error("Error:", err)
		return false
	}
	return d%2 == 1
}

func (r *Receipt) isPurchaseTimeBetweenTwoAndFourNoon(ctx context.Context) bool {
	logger := getLogger(ctx)
	h := r.PurchaseTime[0:2]
	m := r.PurchaseTime[3:5]
	hh, err := strconv.Atoi(h)
	if err != nil {
		logger.Error("Error:", err)
		return false
	}
	mm, err := strconv.Atoi(m)
	if err != nil {
		logger.Error("Error:", err)
		return false
	}
	return (hh >= 14 && hh < 16 && mm > 0 && mm < 60) || (hh == 15 && mm == 0)
}
