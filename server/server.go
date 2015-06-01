package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"

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

	// number of partitions to divide indexes
	Partitions int64 `json:"partitions"`

	// depth of the index tree
	IndexDepth int64 `json:"indexDepth"`

	// maximum payload size in bytes
	PayloadSize int64 `json:"payloadSize"`

	// bucket duration in nano seconds
	// this should be a multiple of `Resolution`
	BucketDuration int64 `json:"bucketDuration"`

	// bucket resolution in nano seconds
	Resolution int64 `json:"resolution"`

	// time duration of a segment
	SegmentDuration int64 `json:"segmentDuration"`

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
	config, err := readConfigFile()
	if err != nil {
		panic(err)
	}

	db, err := kdb.NewDefaultDatabase(kdb.DefaultDatabaseOpts{
		DatabaseName:    config.DatabaseName,
		DataPath:        config.DataPath,
		Partitions:      config.Partitions,
		IndexDepth:      config.IndexDepth,
		PayloadSize:     config.PayloadSize,
		BucketDuration:  config.BucketDuration,
		Resolution:      config.Resolution,
		SegmentDuration: config.SegmentDuration,
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

func readConfigFile() (config *Config, err error) {
	file := flag.String("config", "", "config JSON file")
	flag.Parse()

	if *file == "" {
		return nil, ErrMissingConfigFilePath
	}

	data, err := ioutil.ReadFile(*file)
	if err != nil {
		return nil, err
	}

	config = &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func NewServer(opts ServerOpts) (s *Server) {
	server := bddp.NewServer()
	s = &Server{opts, server}
	server.Method("put", s.handlePut)

	return s
}

func (s *Server) Listen() (err error) {
	log.Print("BDDPServer: listening on ", s.Address)
	return s.server.Listen(s.Address)
}

func (s *Server) handlePut(ctx bddp.MethodContext) {
	defer ctx.SendUpdated()

	params := kmdb.PutRequest_List(*ctx.Params())
	count := params.Len()

	for i := 0; i < count; i++ {
		req := params.At(i)
		pno := req.Partition()
		ts := req.Timestamp()
		vals := req.IndexVals().ToArray()
		pld := req.Payload()

		err := s.Database.Put(ts, pno, vals, pld)
		if err != nil {
			obj := bddp.NewError(ctx.Segment())
			obj.SetError(err.Error())
			ctx.SendError(&obj)
			return
		}
	}

	// TODO: replace below
	obj := ctx.Segment().NewText("")
	ctx.SendResult(&obj)
}
