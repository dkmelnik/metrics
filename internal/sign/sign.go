package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

type Sign struct {
	secret []byte
}

func NewSign(secret string) *Sign {
	return &Sign{[]byte(secret)}
}

func (s *Sign) HashData(data []byte) string {
	h := hmac.New(sha256.New, s.secret)
	h.Write(data)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (s *Sign) Equal(sign string, data []byte) bool {
	calculatedHMAC := s.HashData(data)
	return calculatedHMAC == sign
}
