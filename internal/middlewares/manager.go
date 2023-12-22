package middlewares

import "github.com/dkmelnik/metrics/internal/logger"

type Manager struct {
	logger logger.Logger
}

func NewMiddlewareManager(l logger.Logger) *Manager {
	return &Manager{l}
}
