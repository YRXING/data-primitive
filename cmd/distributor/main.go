package main

import (
	"github.com/YRXING/data-primitive/pkg/distributor"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	d := distributor.NewDistributor()
	d.Run()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT,syscall.SIGTERM)
	<-ch
}
