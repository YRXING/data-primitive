package constants

const (
	// It refers to the prepayment order, and the buyer can pledge the future cargo right.
	ADVANCE = "advance"

	// It refers to account receivable order, and the supplier can pledge the documents
	// of accounts receivable orders.
	ACCOUNT_RECEIVABLE = "account-receivable-order"

	// It refers to financing warehouse order, and the supplier can pledge the inventory.
	FINACING_WAREHOUSE = "financing-warehouse-order"

	NORMAL = "normal"
)

const (
	SUCCESS = 1
	ERROR   = 0
)

const (
	DISTRIBUTOR_SERVICE = "distributor"

	SUPPLIER_SERVICE = "supplier"

	BANK_SERVICE = "bank"

	LOGISTICS_SERVICE = "logistics"
)
