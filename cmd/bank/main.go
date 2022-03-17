package main

import (
	"github.com/YRXING/data-primitive/pkg/bank"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	b := bank.NewBank()
	b.Run()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
