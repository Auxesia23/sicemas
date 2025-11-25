package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"
)

func Encrypt(plain string) string {
	key, _ := base64.StdEncoding.DecodeString(os.Getenv("AES_256_KEY"))
	// Membangun blok AES-256 dari kunci 32 byte.
	block, _ := aes.NewCipher(key)

	// Blok dibungkus mode GCM sehingga kita memperoleh objek AEAD
	// (Authenticated Encryption with Associated Data).
	gcm, _ := cipher.NewGCM(block)

	// Alokasi 12 byte nonce/IV sesuai standar GCM.
	nonce := make([]byte, gcm.NonceSize())

	// 12 byte diisi nilai acak dari crypto/rand
	// agar setiap enkripsi menghasilkan ciphertext berbeda (randomised encryption).
	_, _ = io.ReadFull(rand.Reader, nonce)

	//Fungsi Seal menjalankan dua pekerjaan sekaligus:
	// a. Enkripsi plaintext dengan algoritma Counter Mode,
	// b. Menghitung tag autentikasi 16 byte menggunakan perkalian di lapangan Galois.
	// Hasilnya adalah byte-array tunggal berbentuk [nonce || ciphertext || tag].
	cipherByte := gcm.Seal(nonce, nonce, []byte(plain), nil)

	// Baris 6 – Array byte di-encode ke string heksadesimal agar mudah disimpan di database
	return hex.EncodeToString(cipherByte)
}

func Decrypt(cipherHex string) string {
	key, _ := base64.StdEncoding.DecodeString(os.Getenv("AES_256_KEY"))
	// Mengubah string heksadesimal menjadi byte-array mentah
	cipherByte, _ := hex.DecodeString(cipherHex)

	// Membuat blok AES-256 dari kunci 32 byte (sama seperti saat enkripsi).
	block, _ := aes.NewCipher(key)

	// Blok dibungkus mode GCM sehingga kita memperoleh objek AEAD
	// AEAD mempunyai metode Open
	gcm, _ := cipher.NewGCM(block)

	// Potong array menjadi dua bagian:
	// nonce = 12 byte pertama.
	// cipherText = sisanya
	nonce, cipherText := cipherByte[:gcm.NonceSize()], cipherByte[gcm.NonceSize():]

	// Dekripsi nonce dan cyper text
	plainByte, _ := gcm.Open(nil, nonce, cipherText, nil)

	// Kembalikan byte sebagai string agar bisa dibaca manusia
	return string(plainByte)
}
