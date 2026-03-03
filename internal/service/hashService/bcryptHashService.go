package hashservice

import "golang.org/x/crypto/bcrypt"

type bcryptHashService struct{}

type HashService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) bool
}

// methods

func (s *bcryptHashService) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

func (s *bcryptHashService) VerifyPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
