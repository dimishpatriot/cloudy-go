package runners

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/dimishpatriot/cloudy-go/pkg/services"
	"github.com/dimishpatriot/cloudy-go/pkg/templates/stable"
)

func CircuitBreaker() {
	br := stable.CircuitBreaker(services.RandomSuccess, 2)

	for i := 0; i < 100; i++ {
		fmt.Printf("%d:\t", i)
		res, err := br(context.Background())
		if err != nil {
			fmt.Printf("error:\t\t%s\n", err)
		} else {
			fmt.Printf("success:\t%s\n", res)
		}
		<-time.After(100 * time.Millisecond)
	}
}

func DebounceFirst() {
	dbn := stable.DebounceFirst(services.GetTime, time.Millisecond*200)

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
	dbn := stable.DebounceLast(services.GetTime, time.Millisecond*500)

	for i := 0; i < 20; i++ {
		res, _ := dbn(context.Background())
		fmt.Printf("[%d]\t%s\n", i, res)
		<-time.After(time.Millisecond * 150)
	}
}

func Retry() {
	rr := stable.Retry(services.FlakyService, 3, time.Millisecond*100)

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
	th := stable.Throttle(services.GetTime, 4, time.Millisecond*200)
	runThreeSeries(th)
}

func ThrottleByTokens() {
	th := stable.ThrottleByTokens(services.GetTime, 4, time.Millisecond*200)
	runThreeSeries(th)
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
