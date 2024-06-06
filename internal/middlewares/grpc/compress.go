package grpc

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Compress Compression Interceptor
func (m *MiddlewareManager) Compress() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		if encoding := md.Get("content-encoding"); len(encoding) > 0 && encoding[0] == "gzip" {
			compressedReq, ok := req.([]byte)
			if !ok {
				return nil, status.Error(codes.InvalidArgument, "request is not a byte slice")
			}

			gz, err := gzip.NewReader(bytes.NewReader(compressedReq))
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			defer gz.Close()

			decompressedReq, err := io.ReadAll(gz)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}

			err = json.Unmarshal(decompressedReq, req)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, "failed to unmarshal decompressed data")
			}
		}

		h, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		if acceptEncoding := md.Get("accept-encoding"); len(acceptEncoding) > 0 && strings.Contains(acceptEncoding[0], "gzip") {
			compressedResp := &bytes.Buffer{}
			gz := gzip.NewWriter(compressedResp)
			err = json.NewEncoder(gz).Encode(h)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			gz.Close()

			return compressedResp.Bytes(), nil
		}

		return h, nil
	}
}
