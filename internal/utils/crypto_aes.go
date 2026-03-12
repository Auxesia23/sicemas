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

func Encrypt(plain string) (string, error) {
	// Jika plain text kosong, kembalikan string kosong tanpa mengenkripsinya.
	if plain == "" {
		return "", nil
	}

	key, err := base64.StdEncoding.DecodeString(os.Getenv("AES_256_KEY"))
	if err != nil {
		return "", err
	}
	// Membangun blok AES-256 dari kunci 32 byte.
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Blok dibungkus mode GCM sehingga kita memperoleh objek AEAD
	// (Authenticated Encryption with Associated Data).
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Alokasi 12 byte nonce/IV sesuai standar GCM.
	nonce := make([]byte, gcm.NonceSize())

	// 12 byte diisi nilai acak dari crypto/rand
	// agar setiap enkripsi menghasilkan ciphertext berbeda (randomised encryption).
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	//Fungsi Seal menjalankan dua pekerjaan sekaligus:
	// a. Enkripsi plaintext dengan algoritma Counter Mode,
	// b. Menghitung tag autentikasi 16 byte menggunakan perkalian di lapangan Galois.
	// Hasilnya adalah byte-array tunggal berbentuk [nonce || ciphertext || tag].
	cipherByte := gcm.Seal(nonce, nonce, []byte(plain), nil)

	// Baris 6 – Array byte di-encode ke string heksadesimal agar mudah disimpan di database
	return hex.EncodeToString(cipherByte), nil
}

func Decrypt(cipherHex string) (string, error) {
	// Jika data di database kosong, langsung kembalikan kosong.
	if cipherHex == "" {
		return "", nil
	}

	key, err := base64.StdEncoding.DecodeString(os.Getenv("AES_256_KEY"))
	if err != nil {
		return "", err
	}
	// Mengubah string heksadesimal menjadi byte-array mentah
	cipherByte, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", err
	}

	// Membuat blok AES-256 dari kunci 32 byte (sama seperti saat enkripsi).
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Blok dibungkus mode GCM sehingga kita memperoleh objek AEAD
	// AEAD mempunyai metode Open
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Potong array menjadi dua bagian:
	// nonce = 12 byte pertama.
	// cipherText = sisanya
	nonce, cipherText := cipherByte[:gcm.NonceSize()], cipherByte[gcm.NonceSize():]

	// Dekripsi nonce dan cyper text
	plainByte, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	// Kembalikan byte sebagai string agar bisa dibaca manusia
	return string(plainByte), nil
}
