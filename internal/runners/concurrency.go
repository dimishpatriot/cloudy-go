package runners

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dimishpatriot/cloudy-go/pkg/services"
	"github.com/dimishpatriot/cloudy-go/pkg/templates/concurrency"
)

func FanIn() {
	src := make([]<-chan int, 0)

	for i := 0; i < 40; i += 10 {
		ch := make(chan int)
		src = append(src, ch)

		go func() {
			defer close(ch)
			for j := 1; j <= 3; j++ {
				ch <- i + j
				time.Sleep(time.Millisecond * 100)
			}
		}()

		dest := concurrency.FanIn[int](src...)
		for d := range dest {
			fmt.Printf("value in dest channel: %d\n", d)
		}
	}
}

func FanOut() {
	src := make(chan int)
	dest := concurrency.FanOut[int](src, 5)

	go func() {
		for i := 0; i <= 12; i++ {
			src <- i
		}
		close(src)
	}()

	var wg sync.WaitGroup
	wg.Add(len(dest))

	for i, ch := range dest {
		go func(i int, d <-chan int) {
			defer wg.Done()

			for v := range d {
				fmt.Printf("#%d got %d\n", i, v)
			}
		}(i, ch)
	}

	wg.Wait()
}

func Future() {
	ctx := context.Background()
	future := services.SlowFunc(ctx)
	res, err := future.Result()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(res)
}
