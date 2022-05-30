package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

//Randomly generate Int
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

//Randomly generate String
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

//Randomly generate Owner
func RandomOwner() string {
	return RandomString(6)
}

//Randomly generate Money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

//Randomly generate Currency
func RandomCurrency() string {
	currencies := []string{"WON", "USD", "EUR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// Randomly generate Email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
