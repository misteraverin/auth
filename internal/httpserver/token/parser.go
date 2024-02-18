package token

type Parser interface {
	EncryptToken() string
	Valid() error
}
