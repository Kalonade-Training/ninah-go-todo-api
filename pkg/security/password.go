package security

import "golang.org/x/crypto/bcrypt"

type PasswordService interface {
	Hash(plain string) (string, error)
	Verify(hash, plain string) bool
}

type bcryptService struct{}

func NewBcryptService() PasswordService { return &bcryptService{} }

func (b *bcryptService) Hash(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(bytes), err
}

func (b *bcryptService) Verify(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
