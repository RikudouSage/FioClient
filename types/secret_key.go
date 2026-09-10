package types

import "fmt"

type SecretKey struct {
	bytes  []byte
	string string
}

func NewStringSecretKey(key string) SecretKey {
	return SecretKey{
		string: key,
	}
}

func NewRawSecretKey(key []byte) SecretKey {
	return SecretKey{
		bytes: key,
	}
}

func (receiver SecretKey) String() string {
	if receiver.string != "" {
		return receiver.string
	}

	if receiver.bytes != nil {
		return fmt.Sprintf(`x'%X'`, receiver.bytes)
	}

	return ""
}

func (receiver SecretKey) IsZero() bool {
	return len(receiver.bytes) == 0 && len(receiver.string) == 0
}
