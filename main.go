package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
)

type Circuit func(context.Context) (string, error)

type App struct {
	name string
	f    func()
}

func main() {
	apps := []App{
		{
			name: "runContextCancelling",
			f:    runContextCancelling,
		},
		{
			name: "runBreaker",
			f:    runBreaker,
		},
		{
			name: "runDebounceFirst",
			f:    runDebounceFirst,
		},
	}

	for i, a := range apps {
		fmt.Printf("[%d] %s\n", i, a.name)
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
