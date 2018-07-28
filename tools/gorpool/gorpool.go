package gorpool


import (
"sync"
	"time"
	"log"
)

//Job function
type Job func()

type worker struct {
	workerPool chan *worker
	jobQueue   chan Job
	stop       chan struct{}
}
type Pool struct {
	dispatcher       *dispatcher
	wg               sync.WaitGroup
	enableWaitForAll bool
	workerNum int
	jobNum int
	workerCount int
	time_start time.Time
}

/**
	当前执行中数量
 */
func (p *Pool) GetPoolWorkCount() int{
	return len(p.dispatcher.workerPool)
}

/**
   当前任务总数
 */
func (p *Pool) GetPoolTotalJobCount() int {
	return len(p.dispatcher.jobQueue)
}

type dispatcher struct {
	workerPool chan *worker
	jobQueue   chan Job
	stop       chan struct{}
}

func newWorker(workerPool chan *worker) *worker {
	return &worker{
		workerPool: workerPool,
		jobQueue:   make(chan Job),
		stop:       make(chan struct{}),
	}
}

//one worker start to work
func (w *worker) start() {
	for {
		w.workerPool <- w
		select {
		case job := <-w.jobQueue:
			job()
		case <-w.stop:
			w.stop <- struct{}{}
			return
		}

	}

}

//Dispatch job to free worker
func (dis *dispatcher) dispatch() {
	for {
		select {
		case job := <- dis.jobQueue:
			worker := <- dis.workerPool
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
	return &dispatcher{workerPool: workerPool, jobQueue: jobQueue, stop: make(chan struct{})}
}
//workerNum is worker number of worker pool,on worker have one goroutine
//
//jobNum is job number of job pool
func NewPool(workerNum, jobNum int) *Pool {
	workers := make(chan *worker, workerNum)
	jobs := make(chan Job, jobNum)

	pool := &Pool{
		dispatcher:       newDispatcher(workers, jobs),
		enableWaitForAll: false,
		workerNum:workerNum,
		jobNum:jobNum,
		wg : sync.WaitGroup{},
	}

	return pool

}

//Add one job to job pool
func (p *Pool) AddJob(job Job) {
	if p.enableWaitForAll {
		p.wg.Add(1)
	}
	p.dispatcher.jobQueue <- func() {
		job()
		if p.enableWaitForAll {
			p.wg.Done()
		}
	}
	//log.Printf("len(p.dispatcher.jobQueue):%d",len(p.dispatcher.jobQueue))
	if len(p.dispatcher.jobQueue) >0{
		if p.workerCount < p.workerNum {
			worker := newWorker(p.dispatcher.workerPool)
			p.workerCount++
			go worker.start()
		}
	}
}

func (p *Pool) WaitForAll() {
	if p.enableWaitForAll {
		for i:=0;i<p.workerNum;i++{
			p.AddJob(func() {})
		}
		p.wg.Wait()
		p1 := time.Now()
		log.Printf("本次执行耗时：%dms",(p1.UnixNano()-p.time_start.UnixNano())/1000/1000)
	}
}

func (p *Pool) StopAll() {
	p.dispatcher.stop <- struct{}{}
	<-p.dispatcher.stop
}

//Enable whether enable WaitForAll
func (p *Pool) EnableWaitForAll(enable bool) *Pool {
	p.enableWaitForAll = enable
	return p
}

//Start worker pool and dispatch
func (p *Pool) Start() *Pool {
	p.time_start = time.Now()
	go p.dispatcher.dispatch()
	return p
}
