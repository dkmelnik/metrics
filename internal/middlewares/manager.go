package middlewares

type MiddlewareManager struct {
	sign Signer
}

func NewMiddlewareManager(s Signer) *MiddlewareManager {
	return &MiddlewareManager{s}
}
