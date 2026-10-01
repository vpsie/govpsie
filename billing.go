package govpsie

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var billingPath = "/apps/v2/billing"

type BillingService interface {
	ListInvoices(ctx context.Context, options *ListOptions) ([]Invoice, error)
	ListPurchaseLog(ctx context.Context, options *ListOptions) ([]PurchaseLog, error)
	ApplyVoucher(ctx context.Context, couponIdentifier string) error
	ListAppliedVouchers(ctx context.Context, options *ListOptions) ([]AppliedVouchers, error)
	ListEstimatedUsages(ctx context.Context, options *ListOptions) ([]EstimatedUsages, error)
}

type billingServiceHandler struct {
	client *Client
}

var _ BillingService = &billingServiceHandler{}

type ListInvoicesRoot struct {
	Error bool      `json:"error"`
	Data  []Invoice `json:"data"`
	Total int       `json:"total "`
}

type Invoice struct {
	ID                         int        `json:"id"`
	UserID                     int        `json:"user_id"`
	Date                       time.Time  `json:"date"`
	SerialNumber               string     `json:"serial_number"`
	Total                      FlexString `json:"total"`
	DiscountID                 *int       `json:"discount_id"`
	TaxID                      *int       `json:"tax_id"`
	DiscountValue              *float32   `json:"discount_value"`
	TaxValue                   *float32   `json:"tax_value"`
	TotalAfterDiscountAndTaxes *float32   `json:"total_after_discount_and_taxes"`
	IsPaid                     int        `json:"is_paid"`
	CreatedOn                  time.Time  `json:"created_on"`
	UpdatedAt                  time.Time  `json:"updated_at"`
	Month                      int        `json:"month"`
	Year                       int        `json:"year"`
	Identifier                 string     `json:"identifier"`
	OldInvoice                 int        `json:"old_invoice"`
	IsCustom                   int        `json:"is_custom"`
	EndPeriod                  *string    `json:"end_period"`
	CustomUserName             *string    `json:"custom_user_name"`
	CustomBillAddress          *string    `json:"custom_bill_address"`
	IsHidden                   int        `json:"is_hidden"`
	TaxPercentage              *string    `json:"tax_percentage"`
	DiscountPercentage         *string    `json:"discount_percentage"`
	InvoiceStatus              int        `json:"invoice_status"`
	Notes                      *string    `json:"notes"`
	DueDate                    *string    `json:"due_date"`
	TransactionID              *int       `json:"transaction_id"`
	PaidValue                  *string    `json:"paid_value"`
	PaymentMethod              *string    `json:"payment_method"`
	CustomFooter               *string    `json:"custom_footer"`
	BankDetails                int        `json:"bank_details"`
	VatPercentage              string     `json:"vat_percentage"`
	Username                   string     `json:"username"`
	TaxName                    *string    `json:"tax_name"`
	TaxType                    *string    `json:"tax_type"`
	DiscountName               *string    `json:"discount_name"`
	DiscountType               *string    `json:"discount_type"`
	ResellerName               *string    `json:"reseller_name"`
	ResellerID                 *string    `json:"reseller_id"`
	IsUserHaveReseller         *int       `json:"is_user_have_reseller"`
	StartingDate               string     `json:"startingDate"`
	ClosingDate                string     `json:"closingDate"`
	StartingBalance            float64    `json:"startingBalance"`
	ClosingBalance             float64    `json:"closingBalance"`
}

type PurchaseLog struct {
	UserID        int         `json:"user_id"`
	ID            int         `json:"id"`
	Authority     string      `json:"authority"`
	TransactionID string      `json:"transaction_id"`
	VatPercentage interface{} `json:"vat_percentage"`
	NetAmount     string      `json:"net_amount"`
	VatValue      interface{} `json:"vat_value"`
	TotalAmount   int         `json:"total_amount"`
	Message       string      `json:"message"`
	Amount        int         `json:"amount"`
	CreatedOn     time.Time   `json:"created_on"`
	LastFour      *string     `json:"last_four"`
	CardType      *string     `json:"card_type"`
	Firstname     string      `json:"firstname"`
	Lastname      string      `json:"lastname"`
	BillAddress   string      `json:"bill_address"`
	BillCountry   string      `json:"bill_country"`
	BillCity      string      `json:"bill_city"`
	BillState     string      `json:"bill_state"`
	BillZip       string      `json:"bill_zip"`
	EntityName    *string     `json:"entity_name"`
}

type ListPurchaseLogRoot struct {
	Error bool          `json:"error"`
	Data  []PurchaseLog `json:"data"`
	Total int           `json:"total"`
}

type EstimatedUsages struct {
	ID             int        `json:"id"`
	ProductID      int        `json:"product_id"`
	EntityType     string     `json:"entity_type"`
	EntityID       int        `json:"entity_id"`
	UserID         int        `json:"user_id"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        *string    `json:"end_date"`
	TypeOfTrigger  string     `json:"type_of_trigger"`
	CreatedOn      time.Time  `json:"created_on"`
	UpdatedAt      time.Time  `json:"updated_at"`
	Identifier     string     `json:"identifier"`
	Quantity       int        `json:"quantity"`
	Unit           string     `json:"unit"`
	Price          FlexString `json:"price"`
	CostValue      FlexString `json:"cost_value"`
	CostValueMonth FlexString `json:"cost_value_month"`
	EntityName     string     `json:"entity_name"`
	Description    string     `json:"description"`
}
type ListEstimatedUsagesRoot struct {
	Error       bool              `json:"error"`
	Total       int               `json:"total"`
	Data        []EstimatedUsages `json:"data"`
	BalanceData struct {
		CurrentBalance      float64     `json:"current_balance"`
		BalanceCharged      int         `json:"balance_charged"`
		MonthlyCharge       FlexString  `json:"monthly_charge"`
		ActualMonthlyCharge FlexString  `json:"actual_monthly_charge"`
		AddedWithCc         int         `json:"added_with_cc"`
		AddedWithCcOrPp     int         `json:"added_with_cc_or_pp"`
		BillCity            interface{} `json:"bill_city"`
		BillCountry         interface{} `json:"bill_country"`
		IsPostPaid          int         `json:"is_post_paid"`
	} `json:"balanceData"`
}

type AppliedVouchers struct {
	ID               int     `json:"id"`
	UserID           int     `json:"user_id"`
	CouponId         int     `json:"coupon_id"`
	CouponIdentifier string  `json:"coupon_identifier"`
	Value            float32 `json:"value"`
	Expires          string  `json:"expires"`
}

type ListAppliedVouchersRoot struct {
	Error bool              `json:"error"`
	Data  []AppliedVouchers `json:"data"`
	Total int               `json:"total"`
}

func (s *billingServiceHandler) ListInvoices(ctx context.Context, options *ListOptions) ([]Invoice, error) {
	path := fmt.Sprintf("%s/invoices", billingPath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var root ListInvoicesRoot
	if err := s.client.Do(ctx, req, &root); err != nil {
		return nil, err
	}

	return root.Data, nil
}

func (s *billingServiceHandler) ListPurchaseLog(ctx context.Context, options *ListOptions) ([]PurchaseLog, error) {
	path := fmt.Sprintf("%s/purchase/logs", billingPath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var root ListPurchaseLogRoot
	if err := s.client.Do(ctx, req, &root); err != nil {
		return nil, err
	}

	return root.Data, nil
}

func (s *billingServiceHandler) ListEstimatedUsages(ctx context.Context, options *ListOptions) ([]EstimatedUsages, error) {
	path := fmt.Sprintf("%s/estimated/usages", billingPath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var root ListEstimatedUsagesRoot
	if err := s.client.Do(ctx, req, &root); err != nil {
		return nil, err
	}

	return root.Data, nil
}

func (s *billingServiceHandler) ListAppliedVouchers(ctx context.Context, options *ListOptions) ([]AppliedVouchers, error) {
	path := fmt.Sprintf("%s/coupons", billingPath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var root ListAppliedVouchersRoot
	if err := s.client.Do(ctx, req, &root); err != nil {
		return nil, err
	}

	return root.Data, nil
}

func (s *billingServiceHandler) ApplyVoucher(ctx context.Context, couponIdentifier string) error {
	path := fmt.Sprintf("%s/coupon/add", billingPath)

	applyReq := struct {
		CouponIdentifier string `json:"couponIdentifier"`
	}{
		CouponIdentifier: couponIdentifier,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, &applyReq)
	if err != nil {
		return err
	}

	return s.client.Do(ctx, req, nil)
}
