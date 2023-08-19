package templates

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dimishpatriot/cloudy-go/pkg/services"
)

// CircuitBreaker (Размыкатель цепи) автоматически отключает
// сервисные функции в ответ на вероятную неисправность, чтобы
// предотвратить более крупные или каскадные отказы, устранить повторяющиеся
// ошибки и обеспечить разумную реакцию на ошибки.
func CircuitBreaker(service services.Effector, failThreshold uint) services.Effector {
	countFails := 0
	lastAttempt := time.Now()
	var m sync.RWMutex

	return func(ctx context.Context) (string, error) {
		m.RLock()
		d := countFails - int(failThreshold)
		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Millisecond * 200 << d)
			if time.Now().Before(shouldRetryAt) {
				m.RUnlock()
				return "", fmt.Errorf("service unreachable, fails - %d, lastAttempt %v", countFails, time.Until(lastAttempt))
			}
		}
		m.RUnlock()

		res, err := service(ctx)
		lastAttempt = time.Now()

		m.Lock()
		defer m.Unlock()

		if err != nil {
			countFails++
			return res, err
		}

		countFails = 0
		return res, nil
	}
}
