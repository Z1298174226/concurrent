package pubsub

import (
	"sync"
	"time"
)

type (
	topicFunc func(v interface{}) bool
	subscribe chan interface{}
)

type Publisher struct {
	m          sync.RWMutex
	buffer     int
	timeout    time.Duration
	subscribes map[subscribe]topicFunc
}

func New(buffer int, timeout time.Duration) *Publisher {
	return &Publisher{
		buffer:     buffer,
		timeout:    timeout,
		subscribes: make(map[subscribe]topicFunc),
	}
}

func (p *Publisher) Subscribe() chan interface{} {
	return p.SubscribeTopic(nil)
}

func (p *Publisher) SubscribeTopic(topic topicFunc) chan interface{} {
	ch := make(chan interface{}, p.buffer)
	p.m.Lock()
	defer p.m.Unlock()
	p.subscribes[ch] = topic
	return ch
}

func (p *Publisher) Evict(ch chan interface{}) {
	p.m.Lock()
	defer p.m.Unlock()
	delete(p.subscribes, ch)
	close(ch)
}

func (p *Publisher) Publish(v interface{}) {
	var wg sync.WaitGroup
	p.m.Lock()
	defer p.m.Unlock()
	for sub, topic := range p.subscribes {
		wg.Add(1)
		go p.sendPublish(sub, v, &wg, topic)
	}
	wg.Wait()
}

func (p *Publisher) sendPublish(sub subscribe, v interface{}, wg *sync.WaitGroup, topic topicFunc) {
	defer wg.Done()
	if topic != nil && !topic(v) {
		return
	}
	select {
	case sub <- v:
	case <-time.After(p.timeout):
	}
}

func (p *Publisher) Close() {
	p.m.Lock()
	defer p.m.Unlock()
	for sub := range p.subscribes {
		delete(p.subscribes, sub)
		close(sub)
	}
}
