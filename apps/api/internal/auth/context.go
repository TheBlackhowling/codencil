package auth

import (
	"context"

	"github.com/TheBlackHowling/codencil/apps/api/internal/models"
)

type contextKey string

const userKey contextKey = "codencil_user"

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func UserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(userKey).(*models.User)
	return user, ok
}
