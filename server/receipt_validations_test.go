package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReceipt_Validate(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Receipt
		args    args
		wantErr bool
		errStr  string
	}{
		{"valid receipt", &Receipt{
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
		}, args{context.Background()}, false, ""},
		{"invalid retailer", &Receipt{"*&(@*#()", "2022-01-01", "13:01", []Item{{"Mountain Dew 12PK", "6.49"}}, "6.49"}, args{context.Background()}, true, "invalid retailer"},
		{"invalid purchase date", &Receipt{"Target", "abc", "13:01", []Item{{"Mountain Dew 12PK", "6.49"}}, "6.49"}, args{context.Background()}, true, "invalid purchase date"},
		{"invalid purchase time", &Receipt{"Target", "2022-01-01", "abc", []Item{{"Mountain Dew 12PK", "6.49"}}, "6.49"}, args{context.Background()}, true, "invalid purchase time"},
		{"invalid total", &Receipt{"Target", "2022-01-01", "13:01", []Item{{"Mountain Dew 12PK", "6.49"}}, "abc"}, args{context.Background()}, true, "invalid total"},
		{"invalid item price", &Receipt{"Target", "2022-01-01", "13:01", []Item{{"Mountain Dew 12PK", "abc"}}, "abc"}, args{context.Background()}, true, "invalid item price"},
		{"total mismatch", &Receipt{"Target", "2022-01-01", "13:01", []Item{{"Mountain Dew 12PK", "6.49"}, {"Emils Cheese Pizza", "12.25"}}, "30.39"}, args{context.Background()}, true, "total mismatch with items total"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Validate(tt.args.ctx); (err != nil) != tt.wantErr {
				if tt.wantErr {
					assert.ErrorContains(t, err, tt.errStr)
				}
				t.Errorf("Receipt.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
