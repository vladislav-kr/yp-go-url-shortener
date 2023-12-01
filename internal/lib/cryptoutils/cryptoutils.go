package cryptoutils

import (
	"crypto/rand"
	"math/big"
)

func GenerateRandomString(size int) (string, error) {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_."
	randChars := make([]byte, size)
	for i := 0; i < size; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		randChars[i] = chars[num.Int64()]
	}

	return string(randChars), nil
}
