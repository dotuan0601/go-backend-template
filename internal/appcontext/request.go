package appcontext

import (
	"context"

	"github.com/dotuan0601/go-backend-template/internal/middleware"
)

func GetRequestID(
	ctx context.Context,
) string {

	value, _ := ctx.Value(
		middleware.RequestIDKey,
	).(string)

	return value
}

func GetPath(
	ctx context.Context,
) string {

	value, _ := ctx.Value(
		middleware.PathKey,
	).(string)

	return value
}

func GetMethod(
	ctx context.Context,
) string {

	value, _ := ctx.Value(
		middleware.MethodKey,
	).(string)

	return value
}
