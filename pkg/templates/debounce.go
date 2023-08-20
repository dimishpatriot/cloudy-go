package templates

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dimishpatriot/cloudy-go/pkg/services"
)

// DebounceFirst - при каждом вызове внешней функции – независимо от результата
// – устанавливается временной интервал.
//
// Любой последующий вызов, выполненный до истечения интервала,
// игнорируется, а вызов, выполненный после интервала, передается внутренней функции.
func DebounceFirst(service services.Effector, d time.Duration) services.Effector {
	var threshold time.Time
	var res string
	var err error
	var m sync.Mutex

	return func(ctx context.Context) (string, error) {
		m.Lock()

		defer func() {
			threshold = time.Now().Add(d)
			m.Unlock()
		}()

		if time.Now().Before(threshold) {
			return res, err
		}

		res, err = service(ctx)
		return fmt.Sprint(res, " - ", threshold), err
	}
}

// DebounceLast - при каждом вызове внешней функции – независимо от результата – устанавливается временной интервал.
//
// Реализация будет ждать завершения серии вызовов,
// прежде чем вызовет внутреннюю функцию
func DebounceLast(service services.Effector, d time.Duration) services.Effector {
	var threshold time.Time = time.Now().Add(d)
	var ticker *time.Ticker
	var res string
	var err error
	var once sync.Once
	var m sync.Mutex

	return func(ctx context.Context) (string, error) {
		m.Lock()
		defer m.Unlock()
		ticker = time.NewTicker(time.Millisecond * 10)

		once.Do(func() {
			go func() {
				defer func() {
					m.Lock()
					ticker.Stop()
					once = sync.Once{}
					m.Unlock()
				}()

				for {
					select {
					case <-ticker.C:
						m.Lock()
						if time.Now().After(threshold) {
							res, err = service(ctx)
							m.Unlock()
							threshold = time.Now().Add(d)
							return
						}
						m.Unlock()
					case <-ctx.Done():
						m.Lock()
						res, err = "", ctx.Err()
						m.Unlock()
						return
					}
				}
			}()
		})

		return res, err
	}
}
