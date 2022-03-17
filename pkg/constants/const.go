package constants

const (
	// It refers to the prepayment order, and the buyer can pledge the future cargo right.
	ADVANCE = "advance"

	// It refers to account receivable order, and the supplier can pledge the documents
	// of accounts receivable orders.
	ACCOUNTRECEIVABLE = "account-receivable-order"

	// It refers to financing warehouse order, and the supplier can pledge the inventory.
	FINACINGWAREHOUSE = "financing-warehouse-order"

	NORMAL = "normal"
)

const (
	SUCCESS = 1
	ERROR   = 0
)
