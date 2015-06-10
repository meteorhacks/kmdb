package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/meteorhacks/kdb/dbase"
	"github.com/meteorhacks/kmdb/kmdb"
)

var (
	ErrMissingConfig = errors.New("config file path is missing")
)

func main() {
	fpath := flag.String("config", "", "configuration file (json)")
	flag.Parse()

	if *fpath == "" {
		panic(ErrMissingConfig)
	}

	data, err := ioutil.ReadFile(*fpath)
	if err != nil {
		panic(err)
	}

	config := &kmdb.Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}

	err = validateConfig(config)
	if err != nil {
		panic(err)
	}

	db, err := dbase.New(dbase.Options{
		DatabaseName:   config.DatabaseName,
		DataPath:       config.DataPath,
		IndexDepth:     config.IndexDepth,
		PayloadSize:    config.PayloadSize,
		BucketDuration: config.BucketDuration,
		Resolution:     config.Resolution,
		SegmentSize:    config.SegmentSize,
	})

	if err != nil {
		panic(err)
	}

	s := kmdb.NewServer(db, config)

	// start pprof server
	go startPPROF(config)

	// finally, start the bddp server on main
	// app will exit if bddp server crashes
	log.Println(s.Listen())
}

// TODO: validate config fields
func validateConfig(config *kmdb.Config) (err error) {
	return nil
}

// Listens on port localhost:6060 for pprof http requests
// If debug mode is on, it will listen on all interfaces
func startPPROF(config *kmdb.Config) {
	addr := "localhost:6060"
	if config.DebugMode {
		addr = ":6060"
	}

	log.Println("PPROF: listening on", addr)
	log.Println(http.ListenAndServe(addr, nil))
}
