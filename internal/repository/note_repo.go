package repository

import (
	"context"
	"errors"
	"notes_backend/internal/model"
	"notes_backend/internal/repository/repoutils"

	"gorm.io/gorm"
)

type NoteRepo struct {
	db *gorm.DB
}

type INoteRepo interface {
	GetByID(ctx context.Context, id uint) (*model.Note, error)
	GetAll(ctx context.Context) (*[]model.Note, error)
	Create(ctx context.Context, note *model.Note) error
	Update(ctx context.Context, id uint, note *model.Note) error
	Delete(ctx context.Context, id uint, note *model.Note) error
}

func NewNoteRepo(db *gorm.DB) INoteRepo {
	return &NoteRepo{db}
}

func (r *NoteRepo) GetByID(ctx context.Context, id uint) (*model.Note, error) {
	db := repoutils.DBFromCtx(ctx, r.db)

	note := &model.Note{}
	if err := db.First(note, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return note, nil
}

func (r *NoteRepo) GetAll(ctx context.Context) (*[]model.Note, error) {
	db := repoutils.DBFromCtx(ctx, r.db)

	notes := &[]model.Note{}

	if err := db.Find(notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *NoteRepo) Create(ctx context.Context, note *model.Note) error {
	db := repoutils.DBFromCtx(ctx, r.db)

	return db.Create(&note).Error
}

func (r *NoteRepo) Update(ctx context.Context, id uint, note *model.Note) error {
	db := repoutils.DBFromCtx(ctx, r.db)

	tx := db.
		Model(new(model.Note)).
		Where("id = ?", id).
		Updates(note)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("entity not found or access denied")
	}

	return nil
}

func (r *NoteRepo) Delete(ctx context.Context, id uint, note *model.Note) error {
	db := repoutils.DBFromCtx(ctx, r.db)

	tx := db.
		Model(note).
		Where("id = ?", id).
		Delete(note)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("entity not found or access denied")
	}

	return nil
}
