package stable

import (
	"context"
	"errors"
	"time"

	"github.com/dimishpatriot/cloudy-go/pkg/services"
)

// Retry (Повтор) учитывает возможный временный характер ошибки в распределенной системе
// и осуществляет повторные попытки выполнить неудачную операцию.
func Retry(circuit services.Effector, retries int, delay time.Duration) services.Effector {
	return func(ctx context.Context) (string, error) {
		for r := 0; ; r++ {
			res, err := circuit(ctx)
			if err != nil {
				if r < retries {
					<-time.After(delay)
					continue
				}
				return "", errors.New("too many retries")
			}

			select {
			case <-ctx.Done():
				return "", ctx.Err()
			default:
				return res, nil
			}
		}
	}
}
