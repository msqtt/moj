package crypt

type Cryptor interface {
	Encrypt(string) (string, error)
	Valid(raw, hashed string) (bool, error)
}
