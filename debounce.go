package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// DebounceFirst - при каждом вызове внешней функции – независимо от результата
// – устанавливается временной интервал.
//
// Любой последующий вызов, выполненный до истечения интервала,
// игнорируется, а вызов, выполненный после интервала, передается внутренней функции.
func DebounceFirst(circuit Circuit, d time.Duration) Circuit {
	var m sync.Mutex
	var threshold time.Time
	var res string
	var err error

	return func(ctx context.Context) (string, error) {
		m.Lock()

		defer func() {
			threshold = time.Now().Add(d)
			m.Unlock()
		}()

		if time.Now().Before(threshold) {
			return res, err
		}

		res, err = circuit(ctx)
		return fmt.Sprint(res, " - ", threshold), err
	}
}

func search(ctx context.Context) (string, error) {
	return time.Now().String(), nil
}

func runDebounceFirst() {
	dbn := DebounceFirst(search, time.Millisecond*200)

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			res, _ := dbn(context.Background())
			fmt.Printf("[%d, %d] res: %s\n", i, j, res)
			<-time.After(time.Millisecond * 100)
		}
		<-time.After(time.Millisecond * 150)
	}
}
