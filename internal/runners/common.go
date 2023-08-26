package runners

import (
	"context"
	"fmt"
	"time"

	"github.com/dimishpatriot/cloudy-go/pkg/templates/common"
)

func ContextCancelling() {
	ctx, cancel := context.WithCancel(context.Background())

	out := make(chan string, 1)

	go func() {
		<-time.After(5 * time.Second)
		cancel()
	}()

	for i := 0; i < 5; i++ {
		err := common.Stream(ctx, out)
		if err != nil {
			fmt.Printf("%d:\ttimeout:\t%s\n", i, err)
		} else {
			fmt.Printf("%d:\tresult:\t%s\n", i, <-out)
		}
	}
}
