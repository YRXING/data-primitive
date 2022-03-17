package main

import (
	"github.com/YRXING/data-primitive/pkg/supplier"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	s := supplier.NewSupplier()
	s.Run()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
