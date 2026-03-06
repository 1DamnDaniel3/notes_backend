package repoutils

import (
	"context"
	"log"
	ctxkeys "notes_backend/internal/repository/ctxKeys"
	"reflect"

	"gorm.io/gorm"
)

func DBFromCtx(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(ctxkeys.TxKey{}).(*gorm.DB); ok {
		return tx
	}
	return db
}

func ApplyTenantFilter[T any](
	ctx context.Context,
	db *gorm.DB,
) *gorm.DB {

	userID, ok := ctx.Value(ctxkeys.UserId).(string)
	if !ok {
		log.Println("Shit")
		return db
	}

	// Проверяем: есть ли поле user_id у модели
	if _, ok := reflect.TypeOf(new(T)).Elem().FieldByName("UserID"); ok {
		return db.Where("user_id = ?", userID)
	}

	return db
}
