package runners

import (
	"context"
	"fmt"
	"time"

	"github.com/dimishpatriot/cloudy-go/pkg/services"
	"github.com/dimishpatriot/cloudy-go/pkg/templates"
)

func ContextCancelling() {
	ctx, cancel := context.WithCancel(context.Background())

	out := make(chan string, 1)

	go func() {
		<-time.After(5 * time.Second)
		cancel()
	}()

	for i := 0; i < 5; i++ {
		err := templates.Stream(ctx, out)
		if err != nil {
			fmt.Printf("%d: timeout: %s\n", i, err)
		} else {
			fmt.Printf("%d: result: %s\n", i, <-out)
		}
	}
}

func CircuitBreaker() {
	br := templates.CircuitBreaker(services.RandomSuccess, 2)

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

func DebounceFirst() {
	dbn := templates.DebounceFirst(services.GetTime, time.Millisecond*200)

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			res, _ := dbn(context.Background())
			fmt.Printf("[%d, %d]\t%s\n", i, j, res)
			<-time.After(time.Millisecond * 100)
		}
		<-time.After(time.Millisecond * 150)
	}
}

func DebounceLast() {
	dbn := templates.DebounceLast(services.GetTime, time.Millisecond*500)

	for i := 0; i < 20; i++ {
		res, _ := dbn(context.Background())
		fmt.Printf("[%d]\t%s\n", i, res)
		<-time.After(time.Millisecond * 150)
	}
}
