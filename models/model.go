package models

import datamock "ohlc-aggregator/data-mock"

var (
	QUEUE_SIZE = 100
)

type EventQueue struct {
	Queue chan datamock.OHLC
}
