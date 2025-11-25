package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateOTP6 mengembalikan string 6 digit (000000-999999)
func GenerateOTP6() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1_000_000))
	return fmt.Sprintf("%06d", n.Int64())
}
