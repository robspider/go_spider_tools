package test_chan

import "sync"

type Job func()

type worker struct {
	workerPool chan *worker
	jobQueue   chan Job
	stop       chan struct{}
}

type dispathcher struct {
	workerPool chan *worker
	jobQueue   chan Job
}

type Pool struct {
	dispatcher       *dispathcher
	wg               *sync.WaitGroup
	enableWaitForAll bool
	workerNum int
	workerCount int
}

//Dispatch job to free worker
func (dis *dispathcher) dispatch() {
	for {
		select {
		case job := <-dis.jobQueue:
			worker := <-dis.workerPool
			worker.jobQueue <- job
		case <-dis.stop:
			for i := 0; i < cap(dis.workerPool); i++ {
				worker := <-dis.workerPool
				worker.stop <- struct{}{}
				<-worker.stop
			}
			dis.stop <- struct{}{}
			return
		}
	}
}
func newDispatcher(workerPool chan *worker, jobQueue chan Job) *dispatcher {
	return &dispathcher{workerPool: workerPool, jobQueue: jobQueue, stop: make(chan struct{})}
}



