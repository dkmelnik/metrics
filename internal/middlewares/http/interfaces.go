package http

type Signer interface {
	Equal(sign string, data []byte) bool
}
