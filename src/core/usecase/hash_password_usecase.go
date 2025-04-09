package usecase

import "golang.org/x/crypto/bcrypt"

type IHashPasswordUseCase interface {
	Hash(password string) (string, error)
	CompareHashAndPassword(hashedPassword, password string) error
}
type HashPasswordUseCase struct {
}

func (h HashPasswordUseCase) CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (h HashPasswordUseCase) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func NewHashPasswordUseCase() IHashPasswordUseCase {
	return &HashPasswordUseCase{}
}
