package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/glycerine/go-capnproto"
	"github.com/meteorhacks/bddp"
	"github.com/meteorhacks/kadira-metric-db"
	"github.com/meteorhacks/kdb"
)

var (
	ErrMissingConfigFilePath = errors.New("config file path is missing")
)

type Config struct {
	// database name. Currently only used with naming files
	// can be useful when supporting multiple Databases
	DatabaseName string `json:"databaseName"`

	// place to store data files
	DataPath string `json:"dataPath"`

	// depth of the index tree
	IndexDepth int64 `json:"indexDepth"`

	// payload size should always be equal to this amount
	PayloadSize int64 `json:"payloadSize"`

	// time duration in nano seconds of a range unit
	// this should be a multiple of `Resolution`
	RangeNanos int64 `json:"rangeNanos"`

	// bucket resolution in nano seconds
	Resolution int64 `json:"resolution"`

	// number of records per segment
	SegmentSize int64 `json:"segmentSize"`

	// address to listen for ddp traffic (host:port)
	BDDPAddress string `json:"bddpAddress"`
}

type ServerOpts struct {
	Address  string
	Database kdb.Database
}

type Server struct {
	ServerOpts
	server bddp.Server
}

func main() {
	fpath := flag.String("config", "", "configuration file (json)")
	flag.Parse()

	if *fpath == "" {
		panic(ErrMissingConfigFilePath)
	}

	data, err := ioutil.ReadFile(*fpath)
	if err != nil {
		panic(err)
	}

	config = &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}

	db, err := kdb.NewDefaultDatabase(kdb.DefaultDatabaseOpts{
		DatabaseName: opts.DatabaseName,
		DataPath:     opts.DataPath,
		IndexDepth:   opts.IndexDepth,
		PayloadSize:  opts.PayloadSize,
		RangeNanos:   opts.RangeNanos,
		Resolution:   opts.Resolution,
		SegmentSize:  opts.SegmentSize,
	})

	if err != nil {
		panic(err)
	}

	s := NewServer(ServerOpts{
		Address:  config.BDDPAddress,
		Database: db,
	})

	err = s.Listen()
	if err != nil {
		panic(err)
	}
}

//    Server
// ------------

func NewServer(opts ServerOpts) (s *Server) {
	server := bddp.NewServer()
	s = &Server{opts, server}

	// method handlers
	server.Method("put", s.handlePut)

	return s
}

func (s *Server) Listen() (err error) {
	log.Print("BDDPServer: listening on ", s.Address)
	return s.server.Listen(s.Address)
}

// Method = "put"
// receives a `PutRequest` list and saves metrics 1 by 1
// on success sends a `PutResult` with `ok` set to true
func (s *Server) handlePut(ctx bddp.MethodContext) {
	defer ctx.SendUpdated()

	params := kmdb.PutRequest_List(*ctx.Params())
	count := params.Len()

	for i := 0; i < count; i++ {
		req := params.At(i)
		ts := req.Timestamp()
		vals := req.IndexVals().ToArray()
		pld := req.Payload()

		if err := s.Database.Put(ts, vals, pld); err != nil {
			fmt.Println("Error: Method(put): ", err)
			obj := bddp.NewError(ctx.Segment())
			obj.SetError(err.Error())
			ctx.SendError(&obj)
			return
		}
	}

	seg := ctx.Segment()
	res := kmdb.NewPutResult(seg)
	res.SetOk(true)

	obj := capn.Object(res)
	ctx.SendResult(&obj)
}
