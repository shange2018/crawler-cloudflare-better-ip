package scheduler

import (
	"crawler/worker"
)

type Scheduler struct {
	requestChan chan worker.Request	// 请求通道
	workerChan chan chan worker.Request	// 工作通道
}

type ReadyNotifier interface {
	WorkerReady(chan worker.Request)
}

func (s *Scheduler) Submit(r worker.Request) {
	s.requestChan <- r
}

func (s *Scheduler) WorkerReady(w chan worker.Request) {
	s.workerChan <- w
}

func (s *Scheduler) GetWorkerChan() chan worker.Request {
	return make(chan worker.Request)	// 新建工作通道
}

func (s *Scheduler) Run() {
	s.requestChan = make(chan worker.Request)
	s.workerChan = make(chan chan worker.Request)
	go func() {
		var requestQ []worker.Request	// 请求队列
		var workerQ []chan worker.Request	// 工作队列
		for {
			var activeRequest worker.Request
			var activeWorker chan worker.Request
			if len(requestQ) > 0 &&
				len(workerQ) > 0 {
				activeWorker = workerQ[0]
				activeRequest = requestQ[0]
			}

			select {
			case r := <- s.requestChan:
				requestQ = append(requestQ, r)
			case w := <- s.workerChan:
				workerQ = append(workerQ, w)
			case activeWorker <- activeRequest:
				workerQ = workerQ[1:]
				requestQ = requestQ[1:]
			}

		}
	}()
}