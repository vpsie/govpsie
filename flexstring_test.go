package govpsie

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestFlexStringUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    FlexString
		wantErr bool
	}{
		{name: "string", in: `"16.38"`, want: "16.38"},
		{name: "integer", in: `0`, want: "0"},
		{name: "decimal keeps its digits", in: `16.384066`, want: "16.384066"},
		{name: "negative", in: `-0.3822`, want: "-0.3822"},
		{name: "null", in: `null`, want: ""},
		{name: "empty string", in: `""`, want: ""},
		{name: "boolean is rejected", in: `true`, wantErr: true},
		{name: "object is rejected", in: `{}`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got FlexString
			err := json.Unmarshal([]byte(tt.in), &got)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

// TestNumericMoneyFields reproduces the API returning monetary fields as JSON
// numbers, which previously failed to decode into string fields.
func TestNumericMoneyFields(t *testing.T) {
	t.Run("invoice total", func(t *testing.T) {
		c := setup(t, func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"error":false,"data":[{"id":1,"total":16.384066}],"total":1}`))
		})

		invoices, err := c.Billing.ListInvoices(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListInvoices: %v", err)
		}
		if len(invoices) != 1 || invoices[0].Total != "16.384066" {
			t.Errorf("invoices = %+v", invoices)
		}
	})

	t.Run("profile monthly charge", func(t *testing.T) {
		c := setup(t, func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"error":false,"data":{"id":7,"monthly_charge":0}}`))
		})

		profile, err := c.Profile.GetProfile(context.Background())
		if err != nil {
			t.Fatalf("GetProfile: %v", err)
		}
		if profile.MonthlyCharge != "0" {
			t.Errorf("MonthlyCharge = %q", profile.MonthlyCharge)
		}
	})

	t.Run("estimated usages balance", func(t *testing.T) {
		c := setup(t, func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"error":false,"total":1,"data":[{"id":1,"price":0.0069,"cost_value":"1.20","cost_value_month":4.97}],` +
				`"balanceData":{"current_balance":-0.38,"monthly_charge":0,"actual_monthly_charge":"12.5"}}`))
		})

		// The balance block is decoded alongside the rows; before the fix its
		// numeric monthly_charge failed the whole call.
		usages, err := c.Billing.ListEstimatedUsages(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListEstimatedUsages: %v", err)
		}
		if len(usages) != 1 {
			t.Fatalf("usages = %+v", usages)
		}
		if u := usages[0]; u.Price != "0.0069" || u.CostValue != "1.20" || u.CostValueMonth != "4.97" {
			t.Errorf("usage = %+v", u)
		}
	})
}
