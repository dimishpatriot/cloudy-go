package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/dimishpatriot/cloudy-go/internal/runners"
)

type App struct {
	name string
	f    func()
}

func main() {
	apps := []App{
		{
			name: "runContextCancelling",
			f:    runners.ContextCancelling,
		},
		{
			name: "runBreaker",
			f:    runners.CircuitBreaker,
		},
		{
			name: "runDebounceFirst",
			f:    runners.DebounceFirst,
		},
		{
			name: "runDebounceLast",
			f:    runners.DebounceLast,
		},
		{
			name: "runRetry",
			f:    runners.Retry,
		},
		{
			name: "runThrottle",
			f:    runners.Throttle,
		},
		{
			name: "runThrottleByTokens",
			f:    runners.ThrottleByTokens,
		},
	}

	for i, a := range apps {
		fmt.Printf("[%d]\t%s\n", i, a.name)
	}

	var answer string
	fmt.Print("> ")
	fmt.Scan(&answer)

	idx, err := strconv.Atoi(answer)
	if err != nil || idx < 0 || idx >= len(apps) {
		log.Fatal("bad answer")
	}

	log.Printf("\nstart app: %s\n", apps[idx].name)
	apps[idx].f()
}
