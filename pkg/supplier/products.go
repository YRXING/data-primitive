package supplier

import "github.com/YRXING/data-primitive/pkg/constants"

func SuccessProducts(sn string) *constants.Products  {
	return &constants.Products{
		SupplierName: sn,
		OrderState: constants.SUCCESS,
		Message: "Get products successful.",
	}
}

func ErrorProducts(sn,message string) *constants.Products {
	return &constants.Products{
		SupplierName: sn,
		OrderState: constants.SUCCESS,
		Message: "Get products failed: "+message,
	}
}
