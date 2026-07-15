package xotel

import (
	"context"
	"errors"
)

type ShutdownFunc func(context.Context) error

func SetupTracerProvider(ctx context.Context, serviceName string, endpoint string) (ShutdownFunc, error) {
	if serviceName == "" {
		return nil, errors.New("service name is required")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return func(context.Context) error { return nil }, nil
}
