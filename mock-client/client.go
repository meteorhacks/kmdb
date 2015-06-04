package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/glycerine/go-capnproto"
	"github.com/meteorhacks/bddp"
	"github.com/meteorhacks/kadira-metric-db"
)

const (
	Address     = "localhost:3000"
	Concurrency = 5
	BatchSize   = 10000
)

var (
	counter int64 = 0
)

func main() {
	for i := 0; i < Concurrency; i++ {
		go StartWorker()
	}

	for {
		time.Sleep(time.Second)
		fmt.Println(counter * BatchSize)
		counter = 0
	}
}

func StartWorker() {
	// create a new bddp client
	c := bddp.NewClient()

	// connect to given address
	if err := c.Connect(Address); err != nil {
		fmt.Println("Error: could not connect to " + Address)
		os.Exit(1)
	}

	for {
		SendMetrics(c)
	}
}

func SendMetrics(c bddp.Client) {
	call := c.NewMethodCall("put")
	seg := call.Segment()

	params := kmdb.NewPutRequestList(seg, BatchSize)
	for i := 0; i < BatchSize; i++ {
		pld := make([]byte, 16, 16)
		req := kmdb.NewPutRequest(seg)
		ts := time.Now().UnixNano()

		vals := seg.NewTextList(4)
		vals.Set(0, "a"+strconv.Itoa(rand.Intn(1000)))
		vals.Set(1, "b"+strconv.Itoa(rand.Intn(20)))
		vals.Set(2, "c"+strconv.Itoa(rand.Intn(5)))
		vals.Set(3, "d"+strconv.Itoa(rand.Intn(10)))

		req.SetPayload(pld)
		req.SetTimestamp(ts)
		req.SetIndexVals(vals)
		params.Set(i, req)
	}

	obj := capn.Object(params)
	if _, err := call.Call(obj); err != nil {
		fmt.Println(err)
	}

	atomic.AddInt64(&counter, 1)
}
