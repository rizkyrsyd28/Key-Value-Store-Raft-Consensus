package util

import (
	"math/rand"
	"time"
)

func RandomElectionTimeout(min, max int) time.Duration {
	randomSeconds := rand.Intn(max-min+1) + min
	return time.Duration(randomSeconds) * time.Second
}
