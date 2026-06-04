// Package scripts: utils
package scripts

import (
	"crypto/rand"
	"errors"
	"math/big"
)

const (
	chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	nums  = "1234567890"
)

func RandomString(n int) (string, error) {
	b := make([]byte, n)

	for i := range b {
		r := make([]byte, 1)
		if _, err := rand.Read(r); err != nil {
			return "", err
		}
		b[i] = chars[int(r[0])%len(chars)]
	}

	return string(b), nil
}

func RandomAplhanumericString(n int) (string, error) {
	b := make([]byte, n)

	for i := range b {
		r := make([]byte, 1)
		if _, err := rand.Read(r); err != nil {
			return "", err
		}
		b[i] = (chars + nums)[int(r[0])%len(chars)]
	}

	return string(b), nil
}

// Vibecoded:

func ChooseRandomElements[T any](arr []T, num int) ([]T, error) {
	if num < 0 {
		return nil, errors.New("num cannot be negative")
	}

	if num > len(arr) {
		return nil, errors.New("num exceeds array length")
	}

	// Copy so we don't modify the caller's slice.
	shuffled := make([]T, len(arr))
	copy(shuffled, arr)

	// Fisher-Yates using crypto/rand.
	for i := len(shuffled) - 1; i > 0; i-- {
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return nil, err
		}

		j := int(jBig.Int64())
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled[:num], nil
}

func RandomNumber(max int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}
