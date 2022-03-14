package supplier

const (
	// It refers to the prepayment order, and the buyer can pledge the future cargo right.
	ADVANCE_ORDER = "advance-order"

	// It refers to account receivable order, and the supplier can pledge the documents
	// of accounts receivable orders.
	ACCOUNTRECEIVABLE_ORDER = "account-receivable-order"

	// It refers to financing warehouse order, and the supplier can pledge the inventory.
	FINACINGWAREHOUSE_ORDER = "financing-warehouse-order"

	NORMAL = "normal"
)

type Order struct {
	// There are three different order types, which represent different financing
	// modes, so as to match different linkage paths.
	OrderType string `json:"order_type"`

	OrderPrice int `json:"order_price"`

	OrderCount int `json:"order_count"`
}

type Capital struct {
	BankName string `json:"bank_name"`
	Num      int    `json:"num"`
}
