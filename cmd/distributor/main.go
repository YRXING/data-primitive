package main

import (
	"github.com/YRXING/data-primitive/pkg/distributor"
	"time"
)

func main() {
	d := distributor.NewDistributor()
	d.Run()

	time.Sleep(3*time.Minute)
}
