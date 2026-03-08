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
	GetAll(ctx context.Context, isPublic *bool) (*[]model.Note, error)
	Create(ctx context.Context, note *model.Note) error
	Update(ctx context.Context, id uint, fields map[string]interface{}) error
	Delete(ctx context.Context, id uint, note *model.Note) error
}

func NewNoteRepo(db *gorm.DB) INoteRepo {
	return &NoteRepo{db}
}

// === GetByID
func (r *NoteRepo) GetByID(ctx context.Context, id uint) (*model.Note, error) {
	db := repoutils.DBFromCtx(ctx, r.db)

	note := &model.Note{}
	if err := db.First(note, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return note, nil
}

// === GetAll
func (r *NoteRepo) GetAll(ctx context.Context, isPublic *bool) (*[]model.Note, error) {
	db := repoutils.DBFromCtx(ctx, r.db)
	db, err := repoutils.ApplyTenantFilter[model.Note](ctx, db)
	if err != nil {
		return nil, err
	}
	if isPublic != nil {
		db = db.Where("is_public = ?", *isPublic)
	}

	notes := &[]model.Note{}

	if err := db.Find(notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

// === Create
func (r *NoteRepo) Create(ctx context.Context, note *model.Note) error {
	db := repoutils.DBFromCtx(ctx, r.db)

	return db.Create(&note).Error
}

// === Update
func (r *NoteRepo) Update(ctx context.Context, id uint, fields map[string]interface{}) error {
	db := repoutils.DBFromCtx(ctx, r.db)

	delete(fields, "id")

	tx := db.
		Model(new(model.Note)).
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
