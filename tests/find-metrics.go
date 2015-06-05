package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/glycerine/go-capnproto"
	"github.com/meteorhacks/bddp"
	"github.com/meteorhacks/kmdb/proto"
)

var (
	address     *string
	concurrency *int
	batchsize   *int
	duration    *int64

	counter    = 0
	counterMtx = &sync.Mutex{}
)

func main() {
	address = flag.String("addr", "localhost:3000", "server address")
	concurrency = flag.Int("c", 1, "number of connections to use")
	batchsize = flag.Int("b", 1, "number requests to send per batch")
	duration = flag.Int64("d", 3600000000000, "max test duration")
	flag.Parse()

	for i := 0; i < *concurrency; i++ {
		go StartWorker()
	}

	for {
		time.Sleep(time.Second)
		fmt.Println(*batchsize * counter)

		counterMtx.Lock()
		counter = 0
		counterMtx.Unlock()
	}
}

func StartWorker() {
	// create a new bddp client
	c := bddp.NewClient()

	// connect to given address
	if err := c.Connect(*address); err != nil {
		fmt.Println("Error: could not connect to " + *address)
		os.Exit(1)
	}

	for {
		FindMetrics(c)
	}
}

func FindMetrics(c bddp.Client) {
	call := c.NewMethodCall("get")
	seg := call.Segment()

	params := proto.NewGetRequestList(seg, *batchsize)
	for i := 0; i < *batchsize; i++ {
		req := proto.NewGetRequest(seg)
		end := time.Now().UnixNano() - rand.Int63n(*duration)
		start := end - int64(time.Hour)

		vals := seg.NewTextList(4)
		vals.Set(0, "a"+strconv.Itoa(rand.Intn(1000)))
		vals.Set(1, "b"+strconv.Itoa(rand.Intn(20)))
		vals.Set(2, "c"+strconv.Itoa(rand.Intn(5)))
		vals.Set(3, "") // this will trigger a Find

		req.SetValues(vals)
		req.SetStart(start)
		req.SetEnd(end)
		params.Set(i, req)
	}

	obj := capn.Object(params)
	if _, err := call.Call(obj); err != nil {
		fmt.Println(err)
	}

	counterMtx.Lock()
	counter++
	counterMtx.Unlock()
}
