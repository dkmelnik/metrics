package collect

type Signer interface {
	HashData(data []byte) string
}
