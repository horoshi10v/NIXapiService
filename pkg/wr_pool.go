package pkg

import (
	"log"
	"sync"
)

type WorkerPool struct {
	Count  int
	Sender chan Restaurant
	Ender  chan bool
}

func NewWorkerPool(count int) *WorkerPool {
	return &WorkerPool{
		Count:  count,
		Sender: make(chan Restaurant, count*2),
		Ender:  make(chan bool),
	}
}

func (p *WorkerPool) Run(wg *sync.WaitGroup, handler func(author Restaurant)) {
	defer wg.Done()
	var restaurant Restaurant
	for {
		select {
		case restaurant = <-p.Sender:
			handler(restaurant)
		case <-p.Ender:
			//fmt.Println(<- p.Sender)
			log.Println("I am finish")
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
