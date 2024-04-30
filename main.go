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
	SCRIPTS    = 1000000
)

type Data struct {
	Value int
}

func main() {

	eq := models.EventQueue{
		Queue: make(chan datamock.OHLC, QUEUE_SIZE),
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go aggregator.Aggregate(&eq.Queue)

	go func() {
		for currTime := range ticker.C {
			_ = currTime
			fmt.Printf("Ticker Fired: %v\n", currTime)
			for i := 0; i < SCRIPTS; i++ {
				ohlc := datamock.Seed(i + 1)
				eq.Queue <- ohlc
			}
		}
	}()

	select {}
}
