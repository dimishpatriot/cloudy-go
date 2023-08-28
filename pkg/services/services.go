package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/dimishpatriot/cloudy-go/pkg/templates/concurrency"
)

type Effector func(context.Context) (string, error)

// Фейковые сервисы для тестирования

// RandomSuccess - возвращает удачный результат с шансом 40:
func RandomSuccess(ctx context.Context) (string, error) {
	chance := rand.Float32()
	if chance > 0.4 {
		return fmt.Sprint(chance), nil
	}
	return "", errors.New("bad chances")
}

// FlakyService - возвращает удачный результат с шансом 0.33 и случайной
// задержкой до 500 мс
func FlakyService(ctx context.Context) (string, error) {
	t := time.Duration(rand.Intn(500)) * time.Millisecond
	<-time.After(t)

	chance := rand.Float32()
	if chance > 0.8 {
		return fmt.Sprintf("%.2f %v", chance, t), nil
	}
	return "", errors.New("i'am flaky")
}

// GetTime - пустышка. Возвращает время вызова.
func GetTime(ctx context.Context) (string, error) {
	return time.Now().String(), nil
}

// SlowOps - медленные операции со случайной задержкой
func SlowOps(ctx context.Context) (string, error) {
	t := time.Duration(rand.Intn(6)) * time.Second
	<-time.After(t)
	return t.String(), ctx.Err()
}

func SlowFunc(ctx context.Context) concurrency.Future {
	resCh := make(chan string)
	errCh := make(chan error)

	go func() {
		select {
		case <-time.After(time.Second * 2):
			resCh <- "I slept for 2 seconds"
			errCh <- nil
		case <-ctx.Done():
			resCh <- ""
			errCh <- ctx.Err()
		}
	}()

	return &concurrency.InnerFuture{ResCh: resCh, ErrCh: errCh}
}
