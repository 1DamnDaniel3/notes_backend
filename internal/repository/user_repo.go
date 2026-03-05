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
	Update(ctx context.Context, id uint, fields map[string]interface{}) error
	Delete(ctx context.Context, id uint, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

func NewUserRepo(db *gorm.DB) IUserRepo {
	return &UserRepo{db}
}

// === GetByEmail
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	db := repoutils.DBFromCtx(ctx, r.db)

	user := &model.User{}
	if err := db.First(user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// === GetByID
func (r *UserRepo) GetByID(ctx context.Context, id uint) (*model.User, error) {
	db := repoutils.DBFromCtx(ctx, r.db)

	user := &model.User{}
	if err := db.First(user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// === Create
func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	db := repoutils.DBFromCtx(ctx, r.db)

	return db.Create(&user).Error
}

// === Update
func (r *UserRepo) Update(ctx context.Context, id uint, fields map[string]interface{}) error {
	db := repoutils.DBFromCtx(ctx, r.db)

	delete(fields, "id")

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

// === Delete
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
