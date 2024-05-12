package util

import (
	"math/rand"
	"time"

	"github.com/tamathecxder/randomail"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandomName() string {
	return RandomString(8)
}

func RandomEmail() string {
	return randomail.GenerateRandomEmails(1)[0]
}

func RandomPhone() string {
	return RandomString(9)
}

func RandomDate() time.Time {
	return time.Now().Add(-time.Duration(RandomInt(0, 1000)) * time.Hour)
}

func RandomMoney() int64 {
	return int64(RandomInt(0, 1000))
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
