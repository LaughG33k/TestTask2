package iternal

import (
	"encoding/hex"
	"math/rand"
	"time"

	"golang.org/x/crypto/sha3"
)

func GenerateRefreshToken(length int) string {

	h := sha3.New256()

	result := string(time.Now().UnixMilli())

	for i := 0; i < length; i++ {
		result += string(symbols[rand.Intn(len(symbols))])
	}

	h.Write([]byte(result))

	return hex.EncodeToString(h.Sum(nil))

}
