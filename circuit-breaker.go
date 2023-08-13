package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// CircuitBreaker (Размыкатель цепи) автоматически отключает
// сервисные функции в ответ на вероятную неисправность, чтобы
// предотвратить более крупные или каскадные отказы, устранить повторяющиеся
// ошибки и обеспечить разумную реакцию на ошибки.
func CircuitBreaker(circuit Circuit, failThreshold uint) Circuit {
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

		res, err := circuit(ctx)
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

func getFromOzon(ctx context.Context) (string, error) {
	chance := rand.Float32()
	if chance > 0.6 {
		return fmt.Sprint(chance), nil
	}
	return "", errors.New("bad chances")
}

func runBreaker() {
	br := CircuitBreaker(getFromOzon, 2)
	for i := 0; i < 100; i++ {
		fmt.Printf("%d: ", i)
		res, err := br(context.Background())
		if err != nil {
			fmt.Printf("error: %s\n", err)
		} else {
			fmt.Printf("success: %s\n", res)
		}
		<-time.After(100 * time.Millisecond)
	}
}
