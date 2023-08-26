package common

import (
	"context"
	"time"

	"github.com/dimishpatriot/cloudy-go/pkg/services"
)

func Stream(ctx context.Context, out chan<- string) error {
	newCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	res, err := services.SlowOps(newCtx)
	if err != nil {
		return err
	}

	for {
		select {
		case out <- res:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
