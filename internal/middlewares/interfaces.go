package middlewares

type Signer interface {
	Equal(sign, data []byte) bool
}
