package pool

import (
	"fmt"
	"log"

	"github.com/panjf2000/ants/v2"
)

type Collector struct {
	Work        chan *Job // receives jobs to send to workers
	pool        *ants.Pool
	workerCount int64
}

func StartDispatcher(workerCount int64) Collector {
	input := make(chan *Job) // channel to recieve work

	pool, err := ants.NewPool(int(workerCount))
	if err != nil {
		panic(fmt.Errorf("start dispatcher error: %v", err))
	}

	collector := Collector{
		Work:        input,
		pool:        pool,
		workerCount: workerCount,
	}

	go collector.loopPop()

	return collector
}

func (c Collector) loopPop() {
	for job := range c.Work {

		if c.pool.Running() >= int(c.workerCount+1) {
			log.Println("[pool] worker count nearly full, running:", c.pool.Running(), "workerCount:", c.workerCount)
		}

		err := c.pool.Submit(func(jb *Job) func() {
			return func() {
				jb.JobFunc(0, jb.Data)
			}
		}(job))
		if err != nil {
			log.Printf("[pool] job submit error: %v", err)
		}
	}
}
