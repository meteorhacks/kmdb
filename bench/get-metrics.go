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
	address     *string
	concurrency *int
	batchsize   *int
	duration    *int64
	indexFind   *bool
	randomize   *bool

	counter    = 0
	counterMtx = &sync.Mutex{}
)

func init() {
	address = flag.String("a", "localhost:3000", "address")
	concurrency = flag.Int("c", 5, "concurrency")
	batchsize = flag.Int("b", 10, "requests per batch")
	duration = flag.Int64("d", 3600000000000, "time range")
	indexFind = flag.Bool("f", true, "trigger a find on db")
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
		b, err := c.GetBatch(*batchsize)
		if err != nil {
			log.Println("GET ERROR:", err)
			continue
		}

		for i := 0; i < *batchsize; i++ {
			end := time.Now().UnixNano()
			start := end - *duration
			vals := make([]string, 4, 4)

			if *randomize {
				vals[0] = "a" + strconv.Itoa(rand.Intn(1000))
				vals[1] = "b" + strconv.Itoa(rand.Intn(20))
				vals[2] = "c" + strconv.Itoa(rand.Intn(5))
			} else {
				vals[0] = "a"
				vals[1] = "b"
				vals[2] = "c"
			}

			if *indexFind {
				vals[3] = ""
			} else if *randomize {
				vals[3] = "d" + strconv.Itoa(rand.Intn(10))
			} else {
				vals[3] = "d"
			}

			b.Set(i, vals, start, end)
		}

		if err = b.Send(); err != nil {
			log.Println("GET ERROR:", err)
			continue
		}

		counterMtx.Lock()
		counter++
		counterMtx.Unlock()
	}
}
