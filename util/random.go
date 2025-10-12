package util

import (
	"math/rand"
	"strings"
	"time"
)

const alplabet = "abcdefghijklmnopqrstuvwxyz"

// this function will be called autommatically whenever the util package is used
func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// Random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alplabet)

	for i := 0; i < n; i++ {
		c := alplabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "INR"}
	n := len(currencies)
	return  currencies[rand.Intn(n)]
}