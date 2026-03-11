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
	GetAllPublic(ctx context.Context, page int64) (*[]GetAllPublicBO, error)
}

func NewNoteQueryService(db *gorm.DB) INoteQueryService {
	return &NoteQueryService{db}
}

// BusinessObjects

type GetAllPublicBO struct {
	model.Note
	UserNickname string
}

func (s *NoteQueryService) GetAllPublic(ctx context.Context, page int64) (*[]GetAllPublicBO, error) {
	const pageSize = 20
	offset := (page - 1) * pageSize

	notes := &[]GetAllPublicBO{}

	if err := s.db.
		Table("notes").
		Select("notes.*, users.nickname as user_nickname").
		Joins("join users on users.id = notes.user_id").
		Where("notes.is_public = ?", true).
		Limit(int(pageSize)).
		Offset(int(offset)).
		Scan(notes).Error; err != nil {
		return nil, err
	}

	return notes, nil
}
