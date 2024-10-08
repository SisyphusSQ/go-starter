package routine

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestGO(t *testing.T) {
	fns := make([]func(), 0, 10)
	for i := 0; i < 1086; i++ {
		p := i
		fn := func() {
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("%d\n", p)
		}

		fns = append(fns, fn)
	}

	Go(10, fns)
}

func TestGOE(t *testing.T) {
	fns := make([]func() error, 0, 10)
	for i := 0; i < 10086; i++ {
		fn := func() error {
			return errors.New("test error")
		}

		fns = append(fns, fn)
	}

	err := GoE(5, fns)
	if err != nil {
		t.Logf("err: %v", err)
	}
}

func TestFor(t *testing.T) {
	for i := 0; i < 10086; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("%d\n", i)
	}
}
