package middlewares

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

type MiddlewareManager struct {
	privateKey *rsa.PrivateKey
	sign       Signer
}

func NewMiddlewareManager(privateKeyPath string, signer Signer) (*MiddlewareManager, error) {
	mm := &MiddlewareManager{sign: signer}

	mm.loadPrivateKey(privateKeyPath)

	return mm, nil
}

func (m *MiddlewareManager) loadPrivateKey(path string) error {
	keyFile, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read private key file: %v", err)
	}

	block, _ := pem.Decode(keyFile)
	if block == nil {
		return errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}
	m.privateKey = privateKey
	return nil
}
