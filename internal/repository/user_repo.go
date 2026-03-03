package repository

import (
	"context"
	"errors"
	"notes_backend/internal/model"
	"notes_backend/internal/repository/repoutils"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

type IUserRepo interface {
	GetByID(ctx context.Context, id uint) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, id any, fields map[string]interface{}) error
	Delete(ctx context.Context, id uint, user *model.User) error
}

func NewUserRepo(db *gorm.DB) IUserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) GetByID(ctx context.Context, id uint) (*model.User, error) {
	db := repoutils.DBFromCtx(ctx, r.db)

	user := &model.User{}
	if err := db.First(user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	db := repoutils.DBFromCtx(ctx, r.db)

	return db.Create(&user).Error
}

func (r *UserRepo) Update(ctx context.Context, id any, fields map[string]interface{}) error {
	db := repoutils.DBFromCtx(ctx, r.db)

	delete(fields, "user_id")

	tx := db.
		Model(new(model.User)).
		Where("id = ?", id).
		Updates(fields)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("entity not found or access denied")
	}

	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id uint, user *model.User) error {
	db := repoutils.DBFromCtx(ctx, r.db)

	tx := db.
		Model(user).
		Where("id = ?", id).
		Delete(user)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("entity not found or access denied")
	}

	return nil
}
