package work

import "sync"

type Worker interface {
	Task()
}

type Pool struct {
	work chan Worker
	wg   sync.WaitGroup
}

func New(maxGoroutine int) *Pool {
	p := Pool{
		work: make(chan Worker),
	}
	p.wg.Add(maxGoroutine)
	for i := 0; i < maxGoroutine; i++ {
		go func() {
			defer p.wg.Done()
			for w := range p.work {
				w.Task()
			}
		}()
	}
	return &p
}

func (p *Pool) Run(w Worker) {
	p.work <- w
}

func (p *Pool) ShutDown() {
	close(p.work)
	p.wg.Wait()
}
