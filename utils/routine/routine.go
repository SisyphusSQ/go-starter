package routine

import (
	"sync"
)

func Go(parallel int, fns []func()) {
	var (
		wg sync.WaitGroup
		p  = make(chan struct{}, parallel)
	)
	defer close(p)

	for _, fn := range fns {
		wg.Add(1)
		go func() {
			p <- struct{}{}
			defer func() {
				<-p
				wg.Done()
			}()

			fn()
		}()
	}
	wg.Wait()
}

func GoE(parallel int, fns []func() error) error {
	var (
		err error
		wg  sync.WaitGroup
		p   = make(chan struct{}, parallel)
	)
	defer close(p)

	for _, fn := range fns {
		wg.Add(1)

		go func() {
			p <- struct{}{}
			defer func() {
				<-p
				wg.Done()
			}()

			if err != nil {
				return
			}

			if errG := fn(); errG != nil {
				err = errG
			}
		}()
	}
	wg.Wait()
	return err
}
