package generator

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const AlphanumericLowerDash = "abcdefghijklmnopqrstuvwxyz0123456789-"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomString generates a random string of length n and input list of characters
func RandomString(n uint8, characters string) string {
	sb := strings.Builder{}
	k := len(characters)
	for i := uint8(0); i < n; i++ {
		c := characters[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomURL generates a random URL string
func RandomURL() string {
	return fmt.Sprintf("https://%s.com", RandomString(20, AlphanumericLowerDash))
}
