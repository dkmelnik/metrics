package collect

// Signer is an interface representing an entity capable of hashing data.
// HashData method takes a byte slice as input and returns a hashed representation of the data as a string.
type Signer interface {
	HashData(data []byte) string
}
