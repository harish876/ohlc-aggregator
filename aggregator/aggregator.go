package aggregator

import (
	"fmt"
	"ohlc-aggregator/clock"
	datamock "ohlc-aggregator/data-mock"
	"sync"
	"time"
)

var (
	ohlc *MapLock
)

type OHLC struct {
	Open  float64
	High  float64
	Close float64
	Low   float64
}

type MapData struct {
	Data                OHLC
	Timestamp           time.Time
	OpenValueTimestamp  time.Time
	CloseValueTimestamp time.Time
}

type MapDataHistory struct {
	Prev *MapData
	Curr *MapData
}

func NewMapDataHistory() *MapDataHistory {
	return &MapDataHistory{
		Prev: nil,
		Curr: nil,
	}
}

type MapLock struct {
	Mtx       sync.Mutex
	ScriptMap map[int]*MapDataHistory
}

func Aggregate(queue *chan datamock.OHLC) {
	durationUntilNextTick := clock.RoundToNearestMinute(time.Minute)
	ticker := time.NewTicker(durationUntilNextTick)
	defer ticker.Stop()

	go func() {
		for currTime := range ticker.C {
			_ = currTime
			//Ticker Reset Details
			durationUntilNextTick = clock.RoundToNearestMinute(time.Minute)
			ticker.Reset(durationUntilNextTick)

			start := time.Now()
			ohlc.Mtx.Lock()
			for key, value := range ohlc.ScriptMap {

				//Reset the map
				delete(ohlc.ScriptMap, key)

				if key == 11536 {
					fmt.Println("------------------------- AGGREGATED DATA -------------------------------")
					if value.Prev != nil {
						fmt.Printf("Timestamp for Previous Minute %v\n", value.Prev.Timestamp)
						fmt.Printf("Open: %f High: %f Low: %f Close: %f\n", value.Prev.Data.Open, value.Prev.Data.High, value.Prev.Data.Low, value.Prev.Data.Close)
					}

					if value.Curr != nil {
						fmt.Printf("Timestamp for Current Minute %v\n", value.Curr.Timestamp)
						fmt.Printf("Open: %f High: %f Low: %f Close: %f\n", value.Curr.Data.Open, value.Curr.Data.High, value.Curr.Data.Low, value.Curr.Data.Close)
					}
					fmt.Println("--------------------------------------------------------------------------")
				}
			}
			ohlc.Mtx.Unlock()
			fmt.Println("Total Time Spent in Locked State", time.Since(start))
		}
	}()

	if ohlc == nil {
		ohlc = &MapLock{
			ScriptMap: make(map[int]*MapDataHistory, 0),
			Mtx:       sync.Mutex{},
		}
	}
	for msg := range *queue {
		ohlc.Mtx.Lock()
		currentServerMinute := clock.RoundTimeToPreviousMinute(time.Now())
		currentLTTMinute := clock.RoundTimeToPreviousMinute(msg.LTT)
		scriptId := msg.ScriptId
		if _, ok := ohlc.ScriptMap[scriptId]; !ok {
			ohlc.ScriptMap[scriptId] = NewMapDataHistory()
		}
		timeDelta := int(currentServerMinute.Sub(currentLTTMinute).Minutes())
		val := ohlc.ScriptMap[scriptId]

		switch timeDelta {
		case 0:
			if val.Curr == nil {
				data := MapData{}
				data.Data.Open = msg.LTP
				data.OpenValueTimestamp = msg.LTT
				data.CloseValueTimestamp = msg.LTT
				data.Data.High = msg.LTP
				data.Data.Low = msg.LTP
				data.Data.Close = msg.LTP
				data.Timestamp = currentLTTMinute
				ohlc.ScriptMap[scriptId].Curr = &data
			} else {
				data := val.Curr
				if msg.LTT.Before(data.OpenValueTimestamp) {
					data.Data.Open = msg.LTP
					data.OpenValueTimestamp = msg.LTT
				}
				if msg.LTT.After(data.CloseValueTimestamp) {
					data.Data.Close = msg.LTP
					data.CloseValueTimestamp = msg.LTT
				}
				data.Data.High = max(data.Data.High, msg.LTP)
				data.Data.Low = min(data.Data.Low, msg.LTP)
				ohlc.ScriptMap[scriptId].Curr = data
			}

		case 1:
			if val.Prev == nil {
				data := MapData{}
				data.Data.Open = msg.LTP
				data.OpenValueTimestamp = msg.LTT
				data.CloseValueTimestamp = msg.LTT
				data.Data.High = msg.LTP
				data.Data.Low = msg.LTP
				data.Data.Close = msg.LTP
				data.Timestamp = currentLTTMinute
				ohlc.ScriptMap[scriptId].Prev = &data
			} else {
				data := val.Prev
				if msg.LTT.Before(data.OpenValueTimestamp) {
					data.Data.Open = msg.LTP
					data.OpenValueTimestamp = msg.LTT
				}
				if msg.LTT.After(data.CloseValueTimestamp) {
					data.Data.Close = msg.LTP
					data.CloseValueTimestamp = msg.LTT
				}
				data.Data.High = max(data.Data.High, msg.LTP)
				data.Data.Low = min(data.Data.Low, msg.LTP)
				ohlc.ScriptMap[scriptId].Prev = data
			}
		}
		ohlc.Mtx.Unlock()

	}
}
