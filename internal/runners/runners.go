package runners

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
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

func Retry() {
	rr := templates.Retry(services.FlakyService, 3, time.Millisecond*100)

	for i := 0; i < 10; i++ {
		fmt.Printf("%d: ", i)
		res, err := rr(context.Background())
		if err != nil {
			fmt.Printf("error:\t%s\n", err)
		} else {
			fmt.Printf("success:\t%s\n", res)
		}
		<-time.After(200 * time.Millisecond)
	}
}

func Throttle() {
	th := templates.Throttle(services.GetTime, 4, time.Millisecond*200)
	runThreeSeries(th)
}

func ThrottleByTokens() {
	th := templates.ThrottleByTokens(services.GetTime, 4, time.Millisecond*200)
	runThreeSeries(th)
}

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

		dest := templates.FanIn[int](src...)
		for d := range dest {
			fmt.Printf("value in dest channel: %d\n", d)
		}
	}
}

func FanOut() {
	src := make(chan int)
	dest := templates.FanOut[int](src, 5)

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

func runThreeSeries(th services.Effector) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 30; j++ {
			res, err := th(context.Background())
			t := time.Duration(rand.Intn(150)) * time.Millisecond
			fmt.Printf("[%d %d]\t%d\tres:%s\terr:%v\n", i, j, t.Milliseconds(), res, err)
			<-time.After(t)
		}
		fmt.Println("pause 300")
		<-time.After(time.Millisecond * 300)
	}
}
