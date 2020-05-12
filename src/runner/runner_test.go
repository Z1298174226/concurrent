package runner

import (
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

const timeout = 3 * time.Second

var wg sync.WaitGroup

func TestRun(t *testing.T) {

	wg.Add(3)

	log.Println("Starting Work")

	r := New(timeout)

	r.Add(creatTask(), creatTask(), creatTask())

	if err := FacCtr(&r); err != nil {
		switch err {
		case ErrInterrupt:
			log.Println("Terminate due to interrupt")
			os.Exit(1)
		case ErrTimeout:
			log.Println("Terminate due to timeout")
			os.Exit(2)
		}
	}
	log.Println("Process ended")
	wg.Wait()

	//http.ListenAndServe("0.0.0.0:6060", nil)
}

func creatTask() func(int) {
	return func(id int) {
		defer wg.Done()
		log.Printf("Process - Task #%d .", id)
		time.Sleep(time.Duration(id) * time.Second)
	}
}
