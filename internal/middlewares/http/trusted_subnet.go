package http

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func (m *MiddlewareManager) TrustedSubnet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.trustedSubnet == "" {
			next.ServeHTTP(w, r)
			return
		}

		resolveIP, err := m.resolveIP(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		ok, err := m.checkIPInSubnet(resolveIP.String(), m.trustedSubnet)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if !ok {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *MiddlewareManager) resolveIP(r *http.Request) (net.IP, error) {
	ipStr := r.Header.Get("X-Real-IP")

	ip := net.ParseIP(ipStr)
	if ip == nil {
		ips := r.Header.Get("X-Forwarded-For")
		ipStrs := strings.Split(ips, ",")
		ipStr = ipStrs[0]
		ip = net.ParseIP(ipStr)
	}
	if ip == nil {
		return nil, fmt.Errorf("failed parse ip from http header")
	}
	return ip, nil
}

func (m *MiddlewareManager) checkIPInSubnet(ip, subnet string) (bool, error) {
	_, trustedNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return false, err
	}
	clientIP := net.ParseIP(ip)

	return trustedNet.Contains(clientIP), nil
}
