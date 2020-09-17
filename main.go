package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var reqQue Que

func init() {
	reqQue.data = make(map[string][]chan string)
	reqQue.in = make(chan string)
	reqQue.out = make(chan string)
	reqQue.waitCh = make(chan WaitObj)
	go reqQue.start()
}

func main() {
	log.Println("")
	http.HandleFunc("/hello", hello)
	http.ListenAndServe(":8090", nil)
}

func hello(w http.ResponseWriter, req *http.Request) {
	defer func() {
		recover()
	}()
	keys, ok := req.URL.Query()["key"]
	if !ok {
		fmt.Fprintf(w, "no key\n")
		return
	}
	key := keys[0]
	waites, ok := req.URL.Query()["wait"]
	it, err := strconv.Atoi(waites[0])
	if it == 0 {
		it = 1
	}
	du, err := time.ParseDuration(waites[0])
	if err != nil {
		du = 1 * time.Second
	}
	log.Println(waites[0], du)
	err = reqQue.Wait(key)
	if err != nil {
		fmt.Fprintf(w, err.Error()+" "+key)
		return
	}
	reqQue.add(key)
	defer reqQue.remove(key)
	if waites[0] == "5s" {
		panic("panivc wqewqe")
	}
	log.Println("Url Param 'key' is: " + string(key))
	<-time.After(du)
	fmt.Fprintf(w, "hello "+key)
}
