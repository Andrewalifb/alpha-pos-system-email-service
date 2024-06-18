package dto

type EmailReceiver struct {
	EmailAddress string `json:"email_address"`
}

type HeaderReceipt struct {
	StoreName           string `json:"store_name"`
	StoreAddress        string `json:"store_address"`
	CashierName         string `json:"cashier_name"`
	ReceiptID           string `json:"receipt_id"`
	TransactionDateTime string `json:"transaction_date"`
}

type Items struct {
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	TotalPrice  float64 `json:"total_price"`
}

type BodyReceipt struct {
	Items []Items `json:"items"`
}

type SummaryReceipt struct {
	SubTotalAmount float64 `json:"sub_total_amount"`
	DiscountAmount float64 `json:"discount_amount"`
	TaxAmount      float64 `json:"tax_amount"`
	TotalAmount    float64 `json:"total_amount"`
	CashAmount     float64 `json:"cash_amount"`
	ChangeAmount   float64 `json:"change_amount"`
}

type DigitalReceipt struct {
	Receiver EmailReceiver  `json:"receipt_receiver"`
	Header   HeaderReceipt  `json:"receipt_header"`
	Body     BodyReceipt    `json:"receipt_body"`
	Summary  SummaryReceipt `json:"receipt_summary"`
}
