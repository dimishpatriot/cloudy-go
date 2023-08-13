package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func Stream(ctx context.Context, out chan<- string) error {
	newCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	res, err := SlowOps(newCtx)
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

func SlowOps(ctx context.Context) (string, error) {
	t := time.Duration(rand.Intn(6)) * time.Second
	<-time.After(t)
	return t.String(), ctx.Err()
}

func runContextCancelling() {
	ctx, cancel := context.WithCancel(context.Background())

	out := make(chan string, 1)

	go func() {
		<-time.After(5 * time.Second)
		cancel()
	}()

	for i := 0; i < 5; i++ {
		err := Stream(ctx, out)
		if err != nil {
			fmt.Printf("%d: timeout: %s\n", i, err)
		} else {
			fmt.Printf("%d: result: %s\n", i, <-out)
		}
	}
}
