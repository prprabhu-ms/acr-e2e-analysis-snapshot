package workerpool

import (
	"sync"
)

func Run(poolSize int, actions []func()) {
	q := make(chan poolItem, poolSize-1)

	cleanup := sync.WaitGroup{}
	for i := 0; i < poolSize; i++ {
		cleanup.Add(1)
		go func() {
			defer cleanup.Done()
			for w := range q {
				if w.Done {
					return
				}
				w.Action()
			}

		}()
	}

	for _, a := range actions {
		q <- poolItem{
			Action: a,
			Done:   false,
		}
	}
	for i := 0; i < poolSize; i++ {
		q <- poolItem{Done: true}
	}

	cleanup.Wait()
}

type poolItem struct {
	Action func()
	Done   bool
}
