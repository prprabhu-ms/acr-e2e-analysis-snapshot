package workerpool

import (
	"log"
	"time"
)

func RunBatched(poolSize int, batchSize int, actions []func()) {
	if batchSize <= 0 {
		Run(poolSize, actions)
		return
	}

	for i := 0; i < len(actions); i += batchSize {
		j := i + batchSize
		if j > len(actions) {
			j = len(actions)
		}
		Run(poolSize, actions[i:j])
		log.Printf("Processed %d of %d actions. Waiting for %s before next batch.", j, len(actions), batchWaitTime)
		time.Sleep(batchWaitTime)
	}
}
