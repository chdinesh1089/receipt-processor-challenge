package server

import (
	"context"
	"errors"
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

var regexDollarAmount, _ = regexp.Compile(`^\d+\.\d{2}$`)
var regexRetailer, _ = regexp.Compile(`^[\w\s\-&]+$`)

func (r *Receipt) Validate(ctx context.Context) error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}

	if len(r.Items) == 0 {
		return errors.New("no items - at least one item is required")
	}

	if err := r.validateRetailer(ctx); err != nil {
		return err
	}
	if err := r.validatePurchaseDate(ctx); err != nil {
		return err
	}
	if err := r.validatePurchaseTime(ctx); err != nil {
		return err
	}
	if err := r.validateTotal(ctx); err != nil {
		return err
	}
	return nil
}

func (r *Receipt) validateRetailer(ctx context.Context) error {
	if !regexRetailer.MatchString(r.Retailer) {
		getLogger(ctx).Error("validateRetailer> Invalid retailer:", r.Retailer)
		return errors.New("invalid retailer")
	}
	return nil
}

func (r *Receipt) validatePurchaseDate(ctx context.Context) error {
	_, err := time.Parse("2006-01-02", r.PurchaseDate)
	if err != nil {
		getLogger(ctx).Error("validatePurchaseDate> Error:", err)
		return errors.New("invalid purchase date")
	}
	return nil
}

func (r *Receipt) validatePurchaseTime(ctx context.Context) error {
	_, err := time.Parse("15:04", r.PurchaseTime)
	if err != nil {
		getLogger(ctx).Error("validatePurchaseTime> Error:", err)
		return errors.New("invalid purchase time")
	}
	return nil
}

func (r *Receipt) validateTotal(ctx context.Context) error {
	if !regexDollarAmount.MatchString(r.Total) {
		return errors.New("invalid total")
	}
	total, err := strconv.ParseFloat(r.Total, 64)
	if err != nil {
		getLogger(ctx).Error("validateTotal> Error parsing total:", err)
		return errors.New("invalid total")
	}

	totalCalculated := 0.0
	for _, item := range r.Items {
		if !regexDollarAmount.MatchString(item.Price) {
			return errors.New("invalid item price")
		}
		itemPrice, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			getLogger(ctx).Error("validateItems> Error parsing item price:", err)
			return errors.New("invalid item price")
		}
		totalCalculated += itemPrice
	}

	if math.Abs(total-totalCalculated) > 0.01 {
		getLogger(ctx).Errorf("validateTotal> Total (%v) does not match calculated total (%v)", total, totalCalculated)
		return errors.New("total mismatch with items total")
	}

	return nil
}
