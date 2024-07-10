package crypt

type Cryptor interface {
	Encrypt(string) string
	Valid(raw, hashed string) (bool, error)
}
