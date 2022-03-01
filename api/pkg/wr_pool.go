package pkg

import (
	"NIXSwag/api/internal/models"
	"sync"
)

type WorkerPool struct {
	Count  int
	Sender chan models.Restaurant
	Ender  chan bool
}

func NewWorkerPool(count int) *WorkerPool {
	return &WorkerPool{
		Count:  count,
		Sender: make(chan models.Restaurant, count*2),
		Ender:  make(chan bool),
	}
}

func (p *WorkerPool) Run(wg *sync.WaitGroup, handler func(author models.Restaurant)) {
	defer wg.Done()
	var restaurant models.Restaurant
	for {
		select {
		case restaurant = <-p.Sender:
			handler(restaurant)
		case <-p.Ender:
			//myLog.InfoFile("Routine complete")
			return
		}
	}
}

func (p *WorkerPool) Stop() {
	for i := 0; i < p.Count; i++ {
		p.Ender <- false
	}
	close(p.Sender)
	close(p.Ender)
}
