package templates

import (
	"context"
	"sync"
	"time"

	"github.com/dimishpatriot/cloudy-go/pkg/services"
)

// DebounceLast - при каждом вызове внешней функции – независимо от результата – устанавливается временной интервал.
//
// Реализация будет ждать завершения серии вызовов,
// прежде чем вызовет внутреннюю функцию
func DebounceLast(circuit services.Circuit, d time.Duration) services.Circuit {
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
							res, err = circuit(ctx)
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
