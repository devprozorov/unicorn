package security

type Security struct {
	Tokens *TokenService
	Crypt  *AESCrypt
}

func NewSecurity(jwtSecret, totpEncKeyB64 string) (*Security, error) {
	crypt, err := NewAESCrypt(totpEncKeyB64)
	if err != nil {
		return nil, err
	}
	return &Security{
		Tokens: NewTokenService(jwtSecret),
		Crypt:  crypt,
	}, nil
}
