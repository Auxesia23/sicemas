package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
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

func HashPassword(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ComparePassword(hash, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}
