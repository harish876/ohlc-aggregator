package datamock

import (
	"math/rand"
	"time"
)

type OHLC struct {
	ScriptId  int
	LTP       float64
	LTT       time.Time
	Timestamp time.Time
}

func Seed(scriptId int) OHLC {
	currTime := time.Now()
	return OHLC{
		ScriptId: scriptId,
		LTP:      float64(rand.Intn(100) + 100),
		LTT:      currTime.Add(-1 * time.Duration(rand.Intn(60)) * time.Second),
		//LTT:       currTime,
		Timestamp: currTime,
	}
}
