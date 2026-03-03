package repoutils

import (
	"context"
	ctxkeys "notes_backend/internal/repository/ctxKeys"

	"gorm.io/gorm"
)

func DBFromCtx(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(ctxkeys.TxKey{}).(*gorm.DB); ok {
		return tx
	}
	return db
}
