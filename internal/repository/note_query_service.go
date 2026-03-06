package repository

import (
	"context"
	"notes_backend/internal/model"

	"gorm.io/gorm"
)

type NoteQueryService struct {
	db *gorm.DB
}

type INoteQueryService interface {
	GetAllPublic(ctx context.Context, page int64) (*[]model.Note, error)
}

func NewNoteQueryService(db *gorm.DB) INoteQueryService {
	return &NoteQueryService{db}
}

func (s *NoteQueryService) GetAllPublic(ctx context.Context, page int64) (*[]model.Note, error) {
	const pageSize = 20
	offset := (page - 1) * pageSize

	notes := &[]model.Note{}
	if err := s.db.
		Where("is_public = ?", true).
		Limit(int(pageSize)).
		Offset(int(offset)).
		Find(notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}
