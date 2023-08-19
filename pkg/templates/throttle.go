package templates

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dimishpatriot/cloudy-go/pkg/services"
)

// Throttle - дроссельная заслонка. Ограничивает число запросов
// за единицу времени. Реализация предусматривает скользящее
// окно отметок времени
func Throttle(s services.Effector, max uint, d time.Duration) services.Effector {
	call := make([]time.Time, max)
	var idx uint = 0
	var res string

	return func(ctx context.Context) (string, error) {
		now := time.Now()
		lastIdx := max - 1

		if idx == 0 || now.After(call[idx-1].Add(d)) {
			call[0] = now
			idx = 1
			res, _ = s(ctx)
			return fmt.Sprintf("set 0\t - %s", res), nil
		}

		if idx > lastIdx {
			for i := lastIdx; i > 0; i-- {
				if now.After(call[i].Add(d)) {
					copy(call, call[i+1:])
					call[lastIdx-i] = now
					idx = max - i
					res, _ = s(ctx)
					return fmt.Sprintf("shift %d\t - %s", idx-1, res), nil
				}
			}

			copy(call, call[1:])
			call[max-1] = now
			return "*****", errors.New("429")
		}

		call[idx] = now
		idx++
		res, _ = s(ctx)
		return fmt.Sprintf("add %d\t - %s", idx-1, res), nil
	}
}