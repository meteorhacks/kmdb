package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/meteorhacks/kmdb/kmdb"
)

var (
	dbname      *string
	address     *string
	concurrency *int
	batchsize   *int
	duration    *int64
	indexFind   *bool
	groupData   *bool
	randomize   *bool

	groupBy    = []bool{true, true, true, true}
	counter    = 0
	counterMtx = &sync.Mutex{}
)

func init() {
	dbname = flag.String("n", "test", "database name")
	address = flag.String("a", "localhost:3000", "address")
	concurrency = flag.Int("c", 5, "concurrency")
	batchsize = flag.Int("b", 10, "requests per batch")
	duration = flag.Int64("d", 3600000000000, "time range")
	indexFind = flag.Bool("f", true, "trigger a find on db")
	groupData = flag.Bool("g", true, "group (marge) last field")
	randomize = flag.Bool("r", true, "randomize fields")
	flag.Parse()
}

func main() {
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
	c := kmdb.NewClient(*address)

	// connect to given address
	if err := c.Connect(); err != nil {
		log.Println("ERROR: could not connect to " + *address)
		os.Exit(1)
	}

	batch := &kmdb.GetReqBatch{}
	batch.Batch = make([]*kmdb.GetReq, *batchsize, *batchsize)

	for i := 0; i < *batchsize; i++ {
		req := &kmdb.GetReq{}
		req.Database = *dbname
		req.Fields = make([]string, 4, 4)
		req.Fields[0] = "a"
		req.Fields[1] = "b"
		req.Fields[2] = "c"
		req.Fields[3] = "d"
		req.GroupBy = make([]bool, 4, 4)
		req.GroupBy[0] = true
		req.GroupBy[1] = true
		req.GroupBy[2] = true
		req.GroupBy[3] = true

		if *indexFind {
			req.Fields[3] = ""
		}

		if *groupData {
			req.GroupBy[3] = false
		}

		batch.Batch[i] = req
	}

	for {
		ts2 := time.Now().UnixNano()
		ts1 := ts2 - *duration

		for i := 0; i < *batchsize; i++ {
			req := batch.Batch[i]
			req.StartTime = ts1
			req.EndTime = ts2

			if *randomize {
				req.Fields[0] = "a" + strconv.Itoa(rand.Intn(1000))
				req.Fields[1] = "b" + strconv.Itoa(rand.Intn(20))
				req.Fields[2] = "c" + strconv.Itoa(rand.Intn(5))

				if *indexFind {
					req.Fields[3] = ""
				} else {
					req.Fields[3] = "d" + strconv.Itoa(rand.Intn(10))
				}
			}
		}

		if _, err := c.GetBatch(batch); err != nil {
			log.Println("GET ERROR:", err)
			continue
		}

		counterMtx.Lock()
		counter++
		counterMtx.Unlock()
	}
}
