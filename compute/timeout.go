package main

import (
	"crypto/rand"
	"log"
	"math/big"
	"time"
)

type Timeout struct {
	period time.Duration
	ticker time.Ticker
}

func (t *Timeout) reset() {
	t.ticker.Stop()
	t.ticker = *time.NewTicker(t.period)
}

func createTimeout(period time.Duration) *Timeout {
	return &Timeout{period, *time.NewTicker(period)}
}

func generateRandomInt(lower int, upper int) int {
	l := int64(lower)
	u := int64(upper)
	max := big.NewInt(u - l)
	r, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatalf("Couldn't generate random int!")
	}
	return int(l + r.Int64())
}

func createRandomTimeout(lower int, upper int, period time.Duration) *Timeout {
	randomInt := generateRandomInt(lower, upper)
	period = time.Duration(randomInt) * period
	return createTimeout(period)
}
