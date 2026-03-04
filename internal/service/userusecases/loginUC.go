package userusecases

import (
	"context"
	"errors"
	"notes_backend/internal/model"
	"notes_backend/internal/repository"
	hashservice "notes_backend/internal/service/hashService"
	"notes_backend/internal/service/jwt"
)

type LoginUC struct {
	userRepo   repository.IUserRepo
	hashServ   hashservice.IBcryptHashService
	jwtService jwt.IJWT
}

type ILoginUC interface {
	Execute(ctx context.Context, user *model.User) (string, error)
}

func NewLoginUC(userRepo repository.IUserRepo,
	hashServ hashservice.IBcryptHashService,
	jwtService jwt.IJWT) ILoginUC {
	return &LoginUC{userRepo, hashServ, jwtService}
}

func (uc *LoginUC) Execute(ctx context.Context, user *model.User) (string, error) {

	userEmail := user.Email
	toVerify := user.PasswordHash

	dbUser, err := uc.userRepo.GetByEmail(ctx, userEmail)
	if err != nil {
		return "", err
	}

	if !uc.hashServ.Verify(dbUser.PasswordHash, toVerify) {
		return "", errors.New("Wrong password")
	}

	*user = *dbUser
	user.PasswordHash = ""

	claims := map[string]interface{}{
		"user_id":  dbUser.ID,
		"email":    dbUser.Email,
		"nickname": dbUser.Nickname,
	}

	token, err := uc.jwtService.Sign(claims)
	if err != nil {
		return "", err
	}

	return token, nil

}
