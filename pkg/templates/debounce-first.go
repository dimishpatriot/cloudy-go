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
