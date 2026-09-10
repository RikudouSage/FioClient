package types

import "fmt"

// SecretKey is a SQLCipher database key represented either as text or as raw
// bytes. Its fields are intentionally private to prevent direct mutation.
type SecretKey struct {
	bytes  []byte
	string string
}

// NewStringSecretKey creates a SecretKey from a textual SQLCipher key.
func NewStringSecretKey(key string) SecretKey {
	return SecretKey{
		string: key,
	}
}

// NewRawSecretKey creates a SecretKey from raw key bytes. SQLCipher receives
// the value as a hexadecimal key literal.
func NewRawSecretKey(key []byte) SecretKey {
	return SecretKey{
		bytes: key,
	}
}

// String returns the SQLCipher key expression. Raw keys are returned using
// SQLCipher's hexadecimal key-literal syntax.
func (receiver SecretKey) String() string {
	if receiver.string != "" {
		return receiver.string
	}

	if receiver.bytes != nil {
		return fmt.Sprintf(`x'%X'`, receiver.bytes)
	}

	return ""
}

// IsZero reports whether the key contains neither text nor raw bytes.
func (receiver SecretKey) IsZero() bool {
	return len(receiver.bytes) == 0 && len(receiver.string) == 0
}
