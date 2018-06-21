package modules

import (
	"math/rand"
	"time"
)

func RandInt(n int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(n)
}
