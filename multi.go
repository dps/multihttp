package main

import (
  "fmt"
  "log"
  "math/rand"
  "net/http"
  "strconv"
  "time"
)

func radioactive(oc chan string) {
	for {
		delay := rand.Intn(5)
		time.Sleep(time.Duration(delay) * time.Second)
		select {
		case oc <- strconv.Itoa(delay):
		default:
		}
	}
}

var bus = make(chan string, 1)

type Listener struct {
	ch chan string
}

var listeners = make(chan Listener, 30)

func pump(oc chan string) {
	for {
		val := <-oc
		read := true
		cnt := 0
		for read {
			read = false
			select {
				case listener := <-listeners:
					listener.ch <- val
					cnt = cnt + 1
					read = true
				default:
			}
		}
		log.Printf("%d clients this time.", cnt)
	}
}

func main() {
	bus := make(chan string)
	go func() {radioactive(bus)}()
	go func() {pump(bus)}()
	http.HandleFunc("/x", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("/x request")
		self := Listener{make(chan string)}
		listeners <- self
		timeout := time.After(time.Duration(6) * time.Second)
		select {
		case <- timeout:
		    fmt.Fprintf(w, "--timeout--")
		case v := <- self.ch:
		    fmt.Fprintf(w, "%s", v)
		}
    })
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
    log.Fatal(http.ListenAndServe(":8080", nil))
}