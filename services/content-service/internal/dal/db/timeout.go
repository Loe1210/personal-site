package db

import (
	"context"

	"github.com/Loe1210/personal-site/internal/xresilience"
)

func withRepositoryTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return xresilience.WithTimeout(ctx, xresilience.DefaultMySQLTimeout)
}
