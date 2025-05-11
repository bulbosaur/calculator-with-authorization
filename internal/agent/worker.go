package agent

import (
	"context"
	"log"
	"sync"
	"time"
)

// Mu - мьютекс в рамках микросервиса данного агента
var Mu sync.Mutex

// Worker изолированное выполняет свою задачу по вычислению
func (a *GRPCAgent) Worker(id int) {
	sem := make(chan struct{}, Workers)
	interval := 1 * time.Second
	ctx := context.Background()

	for {
		sem <- struct{}{}
		Mu.Lock()

		task, err := a.getTask(ctx)

		Mu.Unlock()

		if err != nil {
			log.Printf("worker %d: task receiving error: %v", id, err)
			time.Sleep(interval)
			<-sem
			continue
		}

		result, errorMessage, err := a.executeTask(ctx, task)
		if err != nil {
			if task.ID != 0 {
				log.Printf("Worker %d: execution error task ID-%d: %v", id, task.ID, err)
			} else {
				log.Printf("Worker %d: received invalid task (ID=0): %v", id, err)
			}
			time.Sleep(interval)
			<-sem
			continue
		}

		if task.ID != 0 {
			err = a.sendResult(ctx, task.ID, result, errorMessage)
			if err != nil {
				log.Printf("Worker %d: sending error task ID-%d: %v", id, task.ID, err)
			} else {
				log.Printf("Worker %d: success task ID-%d\nresult: %f", id, task.ID, result)
			}
		}

		<-sem
		time.Sleep(interval)
	}
}
