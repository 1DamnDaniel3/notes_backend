package hashservice

import "golang.org/x/crypto/bcrypt"

type BcryptHashService struct{}

type IBcryptHashService interface {
	Hash(password string) (string, error)
	Verify(hashedPassword, password string) bool
}

func NewbcryptHashService() IBcryptHashService {
	return &BcryptHashService{}
}

// methods

func (s *BcryptHashService) Hash(toHash string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(toHash), bcrypt.DefaultCost)
	return string(hashed), err
}

func (s *BcryptHashService) Verify(hashed, toVerify string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(toVerify)) == nil
}
