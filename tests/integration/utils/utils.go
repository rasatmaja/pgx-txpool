package utils

import (
	"math/rand"
	"time"
)

// RandomDuration will return a random duration between min and max
func RandomDuration(min, max int, unit time.Duration) time.Duration {
	return time.Duration(rand.Intn(max-min)+min) * unit
}
