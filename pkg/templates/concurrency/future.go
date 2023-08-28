package concurrency

import "sync"

type Future interface {
	Result() (string, error)
}

type InnerFuture struct {
	once sync.Once
	wg   sync.WaitGroup

	res   string
	err   error
	ResCh <-chan string
	ErrCh <-chan error
}

func (f *InnerFuture) Result() (string, error) {
	f.once.Do(func() {
		f.wg.Add(1)
		defer f.wg.Done()

		f.res = <-f.ResCh
		f.err = <-f.ErrCh
	})
	f.wg.Wait()

	return f.res, f.err
}
