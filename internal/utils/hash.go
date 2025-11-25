package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"log"
	"os"
)

func HashIndex(plain string) []byte {
	pepper := os.Getenv("PEPPER")
	if pepper == "" {
		log.Fatal("PEPPER unset")
	}
	h := hmac.New(sha256.New, []byte(pepper))
	h.Write([]byte(plain))
	return h.Sum(nil)
}
