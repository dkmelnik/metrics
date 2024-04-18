package sign

import (
	"testing"
)

func TestSign_HashData(t *testing.T) {
	secret := "mysecret"
	data := []byte("somedata")

	sign := NewSign(secret)
	expectedHash := "Y3SYzwgET2sR0xOy1GAVWZOVheZvf0GqhTMwMSnFpcs="

	hash := sign.HashData(data)
	if hash != expectedHash {
		t.Errorf("HashData() = %s, want %s", hash, expectedHash)
	}
}

func TestSign_Equal(t *testing.T) {
	secret := "mysecret"
	data := []byte("somedata")

	sign := NewSign(secret)
	expectedHash := "Y3SYzwgET2sR0xOy1GAVWZOVheZvf0GqhTMwMSnFpcs="
	wrongHash := "wronghash"

	if !sign.Equal(expectedHash, data) {
		t.Errorf("Equal() returned false for correct hash")
	}

	if sign.Equal(wrongHash, data) {
		t.Errorf("Equal() returned true for incorrect hash")
	}
}
