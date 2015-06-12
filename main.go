package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/meteorhacks/kdb"
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

	config := &kmdb.ServerConfig{}
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}

	err = validateConfig(config)
	if err != nil {
		panic(err)
	}

	dbs := map[string]kdb.Database{}

	for name, config := range config.Databases {
		db, err := dbase.New(dbase.Options{
			DatabaseName:   name,
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

		dbs[name] = db
	}

	s := kmdb.NewServer(dbs, config)

	// start pprof server
	go startPPROF(config)

	// finally, start the bddp server on main
	// app will exit if bddp server crashes
	log.Println(s.Listen())
}

// TODO: validate config fields
func validateConfig(config *kmdb.ServerConfig) (err error) {
	return nil
}

// Listens on port localhost:6060 for pprof http requests
// If debug mode is on, it will listen on all interfaces
func startPPROF(config *kmdb.ServerConfig) {
	addr := "localhost:6060"
	if config.RemoteDebug {
		addr = ":6060"
	}

	log.Println("PPROF: listening on", addr)
	log.Println(http.ListenAndServe(addr, nil))
}
