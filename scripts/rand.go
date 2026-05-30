// Package scripts: utils
package scripts

import "crypto/rand"

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
