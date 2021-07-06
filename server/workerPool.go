package server

import (
	"sync"
)

type WorkerPool struct {

	work chan IWorker
	wg sync.WaitGroup
}

func NewWorkerPool(poolSize int) *WorkerPool{

	wp := WorkerPool{
		work: make(chan IWorker),
	}

	wp.wg.Add(poolSize)
	for i := 0; i <poolSize; i++ {
		go func() {
			for w := range wp.work{
				w.task()
			}
			wp.wg.Done()
		}()
	}

	return  &wp
}

func (wp *WorkerPool) AddWork(work IWorker)  {
	wp.work <- work
}

func (wp *WorkerPool) Close()  {
	close(wp.work)
	//wait all goroutine
	wp.wg.Wait()
}

