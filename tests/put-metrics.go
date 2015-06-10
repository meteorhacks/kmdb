package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/glycerine/go-capnproto"
	"github.com/meteorhacks/bddp/client"
	"github.com/meteorhacks/kmdb/proto"
)

var (
	address     *string
	concurrency *int
	batchsize   *int

	counter    = 0
	counterMtx = &sync.Mutex{}
)

func main() {
	address = flag.String("addr", "localhost:3000", "server address")
	concurrency = flag.Int("c", 5, "number of connections to use")
	batchsize = flag.Int("b", 10000, "number requests to send per batch")
	flag.Parse()

	for i := 0; i < *concurrency; i++ {
		go StartWorker()
	}

	for {
		time.Sleep(time.Second)
		log.Printf("%d/s\n", *batchsize*counter)

		counterMtx.Lock()
		counter = 0
		counterMtx.Unlock()
	}
}

func StartWorker() {
	// create a new bddp client
	c := client.New(*address)

	// connect to given address
	if err := c.Connect(); err != nil {
		log.Println("ERROR: could not connect to " + *address)
		os.Exit(1)
	}

	for {
		PutMetrics(c)
	}
}

func PutMetrics(c client.Client) {
	call, err := c.Method("put")
	if err != nil {
		return
	}

	seg := call.Segment()

	params := proto.NewPutRequestList(seg, *batchsize)
	for i := 0; i < *batchsize; i++ {
		pld := make([]byte, 16, 16)
		req := proto.NewPutRequest(seg)
		ts := time.Now().UnixNano()

		vals := seg.NewTextList(4)
		vals.Set(0, "a"+strconv.Itoa(rand.Intn(1000)))
		vals.Set(1, "b"+strconv.Itoa(rand.Intn(20)))
		vals.Set(2, "c"+strconv.Itoa(rand.Intn(5)))
		vals.Set(3, "d"+strconv.Itoa(rand.Intn(10)))

		req.SetPayload(pld)
		req.SetTime(ts)
		req.SetValues(vals)
		params.Set(i, req)
	}

	obj := capn.Object(params)
	if _, err := call.Call(obj); err != nil {
		log.Println("PUT ERROR:", err)
	}

	counterMtx.Lock()
	counter++
	counterMtx.Unlock()
}
