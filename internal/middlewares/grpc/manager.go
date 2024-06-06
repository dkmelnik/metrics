package grpc

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

type Signer interface {
	Equal(sign string, data []byte) bool
}

type MiddlewareManager struct {
	trustedSubnet string
	privateKey    *rsa.PrivateKey
	sign          Signer
}

func NewMiddlewareManager(trustedSubnet string, privateKeyPath string, signer Signer) (*MiddlewareManager, error) {
	mm := &MiddlewareManager{sign: signer, trustedSubnet: trustedSubnet}

	err := mm.loadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, err
	}

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
