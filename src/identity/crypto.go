package identity

import (
	srand "crypto/rand"
	"math/rand"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789._-"

	letterIdxBits = 7 // 7 bits to represent a Letter index
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

var randpool rand.Source

func InitRandomPool() error {
	seed, err := srand.Prime(srand.Reader, 256)
	if err != nil {
		return err
	}
	randpool = rand.NewSource(seed.Int64())
	return nil
}

//GetRandomString generates a random base64 string
func GetRandomString(n int) string {
	b := make([]byte, n)
	l := len(letterBytes)
	// A randpool.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, randpool.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randpool.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < l {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
