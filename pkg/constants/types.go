package constants

// There are three different order types, which represent different financing
// modes, so as to match different linkage paths.
type Order struct {
	OrderType string `json:"order_type"`

	OrderPrice int `json:"order_price"`

	OrderCount int `json:"order_count"`
}

type Products struct {
	SupplierName string `json:"supplier_name"`
	OrderState   int    `json:"order_state"`
	Message string `json:"message"`
}

type Capital struct {
	BankName string `json:"bank_name"`
	Num      int    `json:"num"`
}

// Someone who wants to get loan must write this form.
// The bank will according to this form to consider whether you can get loan.
type Form struct {
	// The type indicates which method you will use to get the loan.
	Type string
	Num  int `json:"num"`
}

type Receipt struct {
	Type string `json:"name"`
	Info string `json:"info"`
}