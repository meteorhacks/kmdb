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
	randomize   *bool

	counter    = 0
	counterMtx = &sync.Mutex{}
)

func init() {
	dbname = flag.String("n", "test", "database name")
	address = flag.String("a", "localhost:3000", "address")
	concurrency = flag.Int("c", 5, "concurrency")
	batchsize = flag.Int("b", 10000, "requests per batch")
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

	batch := make([]*kmdb.PutReq, *batchsize, *batchsize)
	for i := 0; i < *batchsize; i++ {
		req := kmdb.NewPutReq()
		req.Fields = make([]string, 4, 4)
		req.Fields[0] = "a"
		req.Fields[1] = "b"
		req.Fields[2] = "c"
		req.Fields[3] = "d"
		batch[i] = req
	}

	for {
		for i := 0; i < *batchsize; i++ {
			req := batch[i]
			req.Timestamp = time.Now().UnixNano()
			if *randomize {
				req.Fields[0] = "a" + strconv.Itoa(rand.Intn(1000))
				req.Fields[1] = "b" + strconv.Itoa(rand.Intn(20))
				req.Fields[2] = "c" + strconv.Itoa(rand.Intn(5))
				req.Fields[3] = "d" + strconv.Itoa(rand.Intn(10))
			}
		}

		if _, err := c.PutBatch(batch); err != nil {
			log.Println("PUT ERROR:", err)
			continue
		}

		counterMtx.Lock()
		counter++
		counterMtx.Unlock()
	}
}
