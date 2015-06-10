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
	payload = []byte{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16,
	}

	address     *string
	concurrency *int
	batchsize   *int
	randomize   *bool

	counter    = 0
	counterMtx = &sync.Mutex{}
)

func init() {
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

	for {
		b, err := c.PutBatch(*batchsize)
		if err != nil {
			log.Println("PUT ERROR:", err)
			continue
		}

		for i := 0; i < *batchsize; i++ {
			ts := time.Now().UnixNano()
			vals := make([]string, 4, 4)

			if *randomize {
				vals[0] = "a" + strconv.Itoa(rand.Intn(1000))
				vals[1] = "b" + strconv.Itoa(rand.Intn(20))
				vals[2] = "c" + strconv.Itoa(rand.Intn(5))
				vals[3] = "d" + strconv.Itoa(rand.Intn(10))
			} else {
				vals[0] = "a"
				vals[1] = "b"
				vals[2] = "c"
				vals[3] = "d"
			}

			b.Set(i, ts, vals, payload)
		}

		if err = b.Send(); err != nil {
			log.Println("PUT ERROR:", err)
			continue
		}

		counterMtx.Lock()
		counter++
		counterMtx.Unlock()
	}
}
