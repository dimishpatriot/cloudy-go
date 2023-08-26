package templates

import "sync"

// FanIn - мультиплексирует несколько входных каналов в один выходной канал
func FanIn[T interface{}](sources ...<-chan T) <-chan T {
	var wg sync.WaitGroup

	dest := make(chan T)
	wg.Add(len(sources))

	for _, ch := range sources {
		go func(c <-chan T) {
			defer wg.Done()

			for value := range c {
				dest <- value
			}
		}(ch)
	}

	// закрывает выходной канал после исчерпания значений в источниках
	go func() {
		wg.Wait()
		close(dest)
	}()

	return dest
}
