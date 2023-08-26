package concurrency

// FanOut - демультиплексирует (разделяет) входной канал на
// заданное число выходных
func FanOut[T interface{}](source <-chan T, n int) []<-chan T {
	dest := make([]<-chan T, n)

	for i := 0; i < n; i++ {
		ch := make(chan T)
		dest[i] = ch

		go func() {
			defer close(ch)
			for value := range source {
				ch <- value
			}
		}()
	}
	return dest
}
