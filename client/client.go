package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/glycerine/go-capnproto"
	"github.com/meteorhacks/bddp"
	"github.com/meteorhacks/kadira-metric-db"
)

var (
	MetricBatchSize    = 1000
	ErrEchoDoesntMatch = errors.New("strings doesn't match")
)

func main() {
	c := Prepare()
	counter := 0

	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Println(counter)
			counter = 0
		}
	}()

	for i := 0; i < 5; i++ {
		go func() {
			for {
				counter++
				SendMetrics(c)
			}
		}()
	}

	select {}
}

func Prepare() (c bddp.Client) {
	// create a client and connect
	c = bddp.NewClient()
	if err := c.Connect("localhost:3000"); err != nil {
		panic(err)
	}

	return c
}

func SendMetrics(c bddp.Client) (err error) {
	call := c.NewMethodCall("put")
	seg := call.Segment()

	params := kmdb.NewPutRequestList(seg, MetricBatchSize)
	for i := 0; i < MetricBatchSize; i++ {
		pld := make([]byte, 16, 16)
		req := kmdb.NewPutRequest(seg)
		ts := time.Now().UnixNano()

		appId := rand.Intn(1000)
		pno := int64(appId % 4)

		vals := seg.NewTextList(4)
		vals.Set(0, "a"+strconv.Itoa(appId))
		vals.Set(1, "t"+strconv.Itoa(rand.Intn(20)))
		vals.Set(2, "h"+strconv.Itoa(rand.Intn(5)))
		vals.Set(3, "d"+strconv.Itoa(rand.Intn(10)))

		req.SetPartition(pno)
		req.SetPayload(pld)
		req.SetTimestamp(ts)
		req.SetIndexVals(vals)
		params.Set(i, req)
	}

	obj := capn.Object(params)
	if _, err = call.Call(obj); err != nil {
		return err
	}

	return nil
}
