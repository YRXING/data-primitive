package main

import (
	"github.com/YRXING/data-primitive/pkg/supplier"
	"time"
)

func main() {
	s := supplier.NewSupplier()
	s.Run()
	time.Sleep(3*time.Minute)
}
