package grpc

import (
	"context"
	"crypto/rsa"
	"encoding/json"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Decryption Interceptor
func (m *MiddlewareManager) Decryption() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		if m.privateKey == nil {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "metadata is not provided")
		}

		encryptedData := md.Get("encrypted-data")
		if len(encryptedData) == 0 {
			return nil, status.Error(codes.InvalidArgument, "no encrypted data provided")
		}

		decryptedData, err := m.decryptData([]byte(encryptedData[0]))
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}

		err = json.Unmarshal(decryptedData, req)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to unmarshal decrypted data")
		}

		return handler(ctx, req)
	}
}

func (m *MiddlewareManager) decryptData(encryptedData []byte) ([]byte, error) {
	decryptedData, err := rsa.DecryptPKCS1v15(nil, m.privateKey, encryptedData)
	if err != nil {
		return nil, err
	}
	return decryptedData, nil
}
