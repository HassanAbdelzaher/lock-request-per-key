package main

import (
	"errors"
	"log"
	"time"
)

type WaitObj struct {
	Key     string
	Channel chan string
}

type Que struct {
	data   map[string][]chan string
	in     chan string
	out    chan string
	waitCh chan WaitObj
}

func (q *Que) start() {
	log.Printf("starting que looop")
	for {
		//log.Println("new loop")
		select {
		case nw := <-q.in:
			//log.Println("in req")
			_, ok := q.data[nw]
			if !ok {
				q.data[nw] = make([]chan string, 0)
			}
		case ol := <-q.out:
			//log.Println("out req")
			channels, ok := q.data[ol]
			if ok {
				for idx := range channels {
					go func(_ch chan string) {
						defer func() {
							recover()
						}()
						_ch <- ol
					}(channels[idx])
				}
				delete(q.data, ol)
			}
		case wt := <-q.waitCh:
			//log.Println("wait req")
			chs, ok := q.data[wt.Key]
			if chs == nil {
				chs = make([]chan string, 0)
			}
			if ok {
				q.data[wt.Key] = append(chs, wt.Channel)
			}
		}
	}
}

func (q *Que) add(key string) {
	q.in <- key
}

func (q *Que) remove(key string) {
	q.out <- key
}

func (q *Que) Wait(key string) error {
	log.Printf("waiting")
	_, ok := q.data[key]
	if !ok {
		log.Printf("no wait")
		return nil
	}
	ch := make(chan string)
	defer close(ch)
	log.Printf("send wait request")
	q.waitCh <- WaitObj{Key: key, Channel: ch}
	log.Printf("wait loop")
	for {
		select {
		case <-time.After(3 * time.Second):
			log.Printf("timeout")
			return errors.New("timeout")
		case <-ch:
			log.Printf("ready")
			return nil
		}
	}
}
