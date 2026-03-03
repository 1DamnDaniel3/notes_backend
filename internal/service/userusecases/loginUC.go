package userusecases

import (
	"notes_backend/internal/model"
	"notes_backend/internal/repository"
)

type LoginUC struct {
	userRepo repository.IUserRepo
}

type ILoginUC interface {
}

func NewLoginUC(userRepo repository.IUserRepo) ILoginUC {
	return &LoginUC{userRepo}
}

func (uc *LoginUC) Execute(user model.User) {

}
