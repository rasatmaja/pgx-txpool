package utils

import (
	"math/rand/v2"
	"time"
)

// RandomDuration will return a random duration between min and max
func RandomDuration(min, max int, unit time.Duration) time.Duration {
	return time.Duration(rand.IntN(max-min)+min) * unit
}
