package pool

import (
	"io"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	maxGoRoutines   = 25
	pooledResources = 3
)

type dbconnection struct {
	ID int32
}

func (db *dbconnection) Close() error {
	log.Println("Close: Connection", db.ID)
	return nil
}

var idCounter int32

func createDBConnection() (io.Closer, error) {
	id := atomic.AddInt32(&idCounter, 1)
	log.Println("Create : New Connection", id)
	return &dbconnection{
		ID: id,
	}, nil
}

func TestPool(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(maxGoRoutines)

	p, err := New(createDBConnection, pooledResources)
	if err != nil {
		log.Println(err)
	}

	for query := 0; query < maxGoRoutines; query++ {
		go func(q int) {
			performQueries(q, p)
			wg.Done()
		}(query)
	}

	wg.Wait()

	log.Println("Program ended")
	p.Close()
}

func performQueries(query int, p *Pool) {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	conn, err := p.Acquire()
	if err != nil {
		log.Println(err)
		return
	}
	defer p.Release(conn)

	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	log.Printf("QID[%d] CID[%d]\n", query, conn.(*dbconnection).ID)
}
