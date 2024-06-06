package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// TrustedSubnet Trusted Subnet Interceptor
func (m *MiddlewareManager) TrustedSubnet() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		if m.trustedSubnet == "" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.PermissionDenied, "metadata is not provided")
		}

		ips := md.Get("x-real-ip")
		if len(ips) == 0 {
			ips = md.Get("x-forwarded-for")
		}

		if len(ips) == 0 {
			return nil, status.Error(codes.PermissionDenied, "IP address not provided")
		}

		ip := net.ParseIP(ips[0])
		if ip == nil {
			return nil, status.Error(codes.PermissionDenied, "invalid IP address")
		}

		ok, err = m.checkIPInSubnet(ip.String(), m.trustedSubnet)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		if !ok {
			return nil, status.Error(codes.PermissionDenied, "access denied")
		}

		return handler(ctx, req)
	}
}

func (m *MiddlewareManager) checkIPInSubnet(ip, subnet string) (bool, error) {
	_, trustedNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return false, err
	}
	clientIP := net.ParseIP(ip)

	return trustedNet.Contains(clientIP), nil
}
