package main

import (
	"log"
	"net/http"
	"time"
)

//func echo(wr http.ResponseWriter, r *http.Request) {
//	msg, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		wr.Write([]byte("echo error"))
//		return
//	}
//
//	writeLen, err := wr.Write(msg)
//	log.Println(writeLen)
//	if err != nil || writeLen != len(msg) {
//		log.Println(err, "write len:", writeLen)
//	}
//}
//
//func main() {
//	http.HandleFunc("/", echo)
//	err := http.ListenAndServe(":8080", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//}

func hello(wr http.ResponseWriter, r *http.Request) {
	wr.Write([]byte("hello"))
}

func timeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()

		// next handler
		next.ServeHTTP(wr, r)

		timeElapsed := time.Since(timeStart)
		log.Println(timeElapsed)
	})
}

type middleware func(http.Handler) http.Handler

type Router struct {
	middlewareChain []middleware
	mux             map[string]http.Handler
}

func NewRouter() *Router {
	return &Router{
		mux: make(map[string]http.Handler),
	}
}

func (r *Router) Use(m middleware) {
	r.middlewareChain = append(r.middlewareChain, m)
}

func (r *Router) Add(route string, h http.Handler) {
	var mergedHandler = h

	for i := len(r.middlewareChain) - 1; i >= 0; i-- {
		mergedHandler = r.middlewareChain[i](mergedHandler)
	}

	r.mux[route] = mergedHandler
}

func main() {
	r := NewRouter()
	r.Use(timeMiddleware)
	//r.Use(timeout)
	//r.Use(ratelimit)
	r.Add("/", http.HandlerFunc(hello))
	http.Handle("/", r.mux["/"])
	http.ListenAndServe(":8080", nil)

}
