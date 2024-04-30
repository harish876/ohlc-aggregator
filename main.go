package main

import (
	"fmt"
	"ohlc-aggregator/aggregator"
	datamock "ohlc-aggregator/data-mock"
	"ohlc-aggregator/models"
	"time"
)

var (
	QUEUE_SIZE = 1000000
)

type Data struct {
	Value int
}

func main() {
	scriptIds := []int{11536}

	eq := models.EventQueue{
		Queue: make(chan datamock.OHLC, QUEUE_SIZE),
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go aggregator.Aggregate(&eq.Queue)

	go func() {
		for currTime := range ticker.C {
			for _, scriptId := range scriptIds {
				ohlc := datamock.Seed(scriptId)
				fmt.Printf("Ticker: %v , LTP: %f, LTT: %v\n", currTime, ohlc.LTP, ohlc.LTT)
				eq.Queue <- ohlc
			}
		}
	}()

	select {}
}
