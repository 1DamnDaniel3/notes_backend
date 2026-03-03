package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"size:255;unique;not null" json:"email"`
	Nickname     string    `gorm:"size:100;unique;not null" json:"nickname"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Notes        []Note    `gorm:"constraint:OnDelete:CASCADE"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.PasswordHash != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.PasswordHash = string(hashed)
	}
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if tx.Statement.Changed("PasswordHash") && u.PasswordHash != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.PasswordHash = string(hashed)
	}
	return nil
}
