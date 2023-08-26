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
			name: "ContextCancelling",
			f:    runners.ContextCancelling,
		},
		{
			name: "CircuitBreaker",
			f:    runners.CircuitBreaker,
		},
		{
			name: "DebounceFirst",
			f:    runners.DebounceFirst,
		},
		{
			name: "DebounceLast",
			f:    runners.DebounceLast,
		},
		{
			name: "Retry",
			f:    runners.Retry,
		},
		{
			name: "Throttle",
			f:    runners.Throttle,
		},
		{
			name: "ThrottleByTokens",
			f:    runners.ThrottleByTokens,
		},
		{
			name: "FunIn",
			f:    runners.FanIn,
		},
		{
			name: "FunOut",
			f:    runners.FanOut,
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

	log.Printf("\n\nstart template - %s -\n\n", apps[idx].name)
	apps[idx].f()
}
