package hashStorage

import (
	"time"
)

type Counter struct {
	id               int
	timeSeriesValues [][]int
	expirationTime   int
	timer            int
	timePeriod       int
}

func (c *Counter) timerCounter(tmp int) {
	for {
		select {
		case <-time.After(time.Duration(tmp) * time.Second):
			c.timer += 1
			if c.timer >= c.expirationTime {
				c.timer = 0
			}
			c.timeSeriesValues[c.timer] = make([]int, 0, c.timePeriod)
		}
	}
}

func NewCounter(timePer, expiration int) Counter {
	exp := (expiration * 60) / timePer
	tmp := make([][]int, exp)
	for i := range tmp {
		tmp[i] = make([]int, 0, timePer)
	}
	c := Counter{timeSeriesValues: tmp,
		expirationTime: exp,
		timePeriod:     timePer,
	}
	if exp > 0 {
		go c.timerCounter(timePer)
	} else {
		c.expirationTime = -1
		c.timeSeriesValues = append(c.timeSeriesValues, make([]int, 1))
	}
	return c
}

func (c *Counter) IncreasingTheCounterValue(value int) {
	if c.expirationTime > -1 {
		c.timeSeriesValues[c.timer] = append(c.timeSeriesValues[c.timer], value)
		return
	}
	c.timeSeriesValues[0][0] += value

}

func (c *Counter) LenTimeSeries() int {
	return len(c.timeSeriesValues)
}

func (c *Counter) SumActual() int {
	result := 0
	for _, series := range c.timeSeriesValues {
		for _, ser := range series {
			result += ser
		}
	}
	return result
}
