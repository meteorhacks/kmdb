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

	for {
		b, err := c.PutBatch(*batchsize)
		if err != nil {
			log.Println("PUT ERROR:", err)
			continue
		}

		for i := 0; i < *batchsize; i++ {
			ts := time.Now().UnixNano()
			fields := make([]string, 4, 4)

			if *randomize {
				fields[0] = "a" + strconv.Itoa(rand.Intn(1000))
				fields[1] = "b" + strconv.Itoa(rand.Intn(20))
				fields[2] = "c" + strconv.Itoa(rand.Intn(5))
				fields[3] = "d" + strconv.Itoa(rand.Intn(10))
			} else {
				fields[0] = "a"
				fields[1] = "b"
				fields[2] = "c"
				fields[3] = "d"
			}

			b.Set(i, *dbname, ts, fields, 1, 1)
		}

		if _, err = b.Send(); err != nil {
			log.Println("PUT ERROR:", err)
			continue
		}

		counterMtx.Lock()
		counter++
		counterMtx.Unlock()
	}
}
