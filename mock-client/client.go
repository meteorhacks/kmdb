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
	Concurrency = 10
	BatchSize   = 1000

	NumApps  = 1000
	NumTypes = 20
	NumHosts = 5
	NumData  = 10
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
		fmt.Println(counter)
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

func SendMetrics(c bddp.Client) (err error) {
	call := c.NewMethodCall("put")
	seg := call.Segment()

	params := kmdb.NewPutRequestList(seg, BatchSize)
	for i := 0; i < BatchSize; i++ {
		pld := make([]byte, 16, 16)
		req := kmdb.NewPutRequest(seg)
		ts := time.Now().UnixNano()

		vals := seg.NewTextList(4)
		vals.Set(0, "a"+strconv.Itoa(rand.Intn(NumApps)))
		vals.Set(1, "t"+strconv.Itoa(rand.Intn(NumTypes)))
		vals.Set(2, "h"+strconv.Itoa(rand.Intn(NumHosts)))
		vals.Set(3, "d"+strconv.Itoa(rand.Intn(NumData)))

		req.SetPayload(pld)
		req.SetTimestamp(ts)
		req.SetIndexVals(vals)
		params.Set(i, req)
	}

	obj := capn.Object(params)
	if _, err = call.Call(obj); err != nil {
		return err
	}

	atomic.AddInt64(&counter, 1)

	return nil
}
