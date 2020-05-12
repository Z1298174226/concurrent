package pubsub

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestPub(t *testing.T) {
	p := New(10, 100*time.Millisecond)
	defer p.Close()

	var wg sync.WaitGroup

	all := p.Subscribe()
	golang := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "golang")
		}
		return false
	})

	p.Publish("hello,  world!")
	p.Publish("hello, golang!")

	wg.Add(2)
	go func() {
		for msg := range all {
			fmt.Println("all:", msg)
		}
	}()

	go func() {
		for msg := range golang {
			fmt.Println("golang:", msg)
		}
	}()

	time.Sleep(3 * time.Second)

}
