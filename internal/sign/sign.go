package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"hash"
)

type Sign struct {
	hash hash.Hash
}

func NewSign(secret string) *Sign {
	return &Sign{hmac.New(sha256.New, []byte(secret))}
}

func (s *Sign) HashData(data []byte) string {
	return hex.EncodeToString(s.hash.Sum(data))
}

func (s *Sign) Equal(sign, data []byte) bool {
	return subtle.ConstantTimeCompare(sign, s.hash.Sum(data)) == 1
}
