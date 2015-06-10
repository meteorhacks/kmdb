package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"

	"github.com/glycerine/go-capnproto"
	"github.com/meteorhacks/bddp"
	"github.com/meteorhacks/bddp/server"
	"github.com/meteorhacks/kdb"
	"github.com/meteorhacks/kdb/dbase"
	"github.com/meteorhacks/kmdb/proto"

	"net/http"
	_ "net/http/pprof"
)

var (
	ErrMissingConfig = errors.New("config file path is missing")
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
	BucketDuration int64 `json:"bucketDuration"`

	// bucket resolution in nano seconds
	Resolution int64 `json:"resolution"`

	// number of records per segment
	SegmentSize int64 `json:"segmentSize"`

	// enable pprof on ":6060" instead of "localhost:6060".
	DebugMode bool `json:"debugMode"`

	// address to listen for ddp traffic (host:port)
	BDDPAddress string `json:"bddpAddress"`
}

type Server struct {
	Database kdb.Database
	server   server.Server
	config   *Config
}

func main() {
	fpath := flag.String("config", "config.json", "configuration file (json)")
	flag.Parse()

	if *fpath == "" {
		panic(ErrMissingConfig)
	}

	data, err := ioutil.ReadFile(*fpath)
	if err != nil {
		panic(err)
	}

	config := &Config{}
	err = json.Unmarshal(data, config)
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

	s := NewServer(db, config)

	// start a pprof server
	go func() {
		addr := "localhost:6060"
		if config.DebugMode {
			addr = ":6060"
		}

		log.Println("PPROF: listening on", addr)
		log.Println(http.ListenAndServe(addr, nil))
	}()

	// finally, start the bddp server on main
	// app will exit if bddp server crashes
	log.Println(s.Listen())
}

//    Server
// ------------

func NewServer(db kdb.Database, config *Config) (s *Server) {
	srvr := server.New(config.BDDPAddress)
	s = &Server{db, srvr, config}

	// method handlers
	srvr.Method("put", s.handlePut)
	srvr.Method("get", s.handleGet)

	return s
}

func (s *Server) Listen() (err error) {
	log.Println("BDDP:  listening on", s.config.BDDPAddress)
	return s.server.Listen()
}

// Method = "put"
// receives a `PutRequest` list and saves metrics 1 by 1
// on success sends a `PutResult` with `ok` set to true
func (s *Server) handlePut(ctx server.MContext) {
	defer ctx.SendUpdated()

	params := proto.PutRequest_List(*ctx.Params())
	count := params.Len()

	valsCount := int(s.config.IndexDepth)
	vals := make([]string, valsCount, valsCount)

	for i := 0; i < count; i++ {
		req := params.At(i)
		ts := req.Time()
		pld := req.Payload()

		vals := vals[:0]
		for j := 0; j < valsCount; j++ {
			vals = append(vals, req.Values().At(j))
		}

		if err := s.Database.Put(ts, vals, pld); err != nil {
			s.methodError(ctx, err)
			return
		}
	}

	seg := ctx.Segment()
	res := proto.NewPutResult(seg)
	res.SetOk(true)

	obj := capn.Object(res)
	ctx.SendResult(&obj)
}

// Method = "get"
// receives a `GetRequest` list and responds with a list of `GetResult`
// uses either `db.Get()` or `db.Find()` to get data from the database
func (s *Server) handleGet(ctx server.MContext) {
	defer ctx.SendUpdated()

	params := proto.GetRequest_List(*ctx.Params())
	count := params.Len()

	seg := ctx.Segment()
	out := proto.NewGetResultList(seg, count)

	for i := 0; i < count; i++ {
		req := params.At(i)
		start := req.Start()
		end := req.End()
		vals := req.Values().ToArray()

		res := proto.NewGetResult(seg)
		out.Set(i, res)

		var resItems proto.ResultItem_List

		// use the `Get` method only if all values are set
		// otherwise use the more costly `Find` method
		if s.canUseGet(vals) {
			data, err := s.Database.Get(start, end, vals)
			if err != nil {
				continue
			}

			resItems = proto.NewResultItemList(seg, 1)
			resItem := s.makeResultItem(seg, data, vals)
			resItems.Set(0, *resItem)
		} else {
			dataMap, err := s.Database.Find(start, end, vals)
			if err != nil {
				continue
			}

			numItems := len(dataMap)
			resItems = proto.NewResultItemList(seg, numItems)

			counter := 0
			for el, data := range dataMap {
				resItem := s.makeResultItem(seg, data, el.Values)
				resItems.Set(counter, *resItem)
				counter++
			}
		}

		res.SetOk(true)
		res.SetData(resItems)
	}

	obj := capn.Object(out)
	ctx.SendResult(&obj)
}

func (s *Server) canUseGet(vals []string) (can bool) {
	for _, v := range vals {
		if v == "" {
			return false
		}
	}

	return true
}

func (s *Server) makeResultItem(seg *capn.Segment, data [][]byte, vals []string) (resItem *proto.ResultItem) {
	item := proto.NewResultItem(seg)
	resItem = &item

	resData := seg.NewDataList(len(data))
	resItem.SetData(resData)
	for j, d := range data {
		resData.Set(j, d)
	}

	resVals := seg.NewTextList(len(vals))
	resItem.SetValues(resVals)
	for j, v := range vals {
		resVals.Set(j, v)
	}

	return resItem
}

func (s *Server) methodError(ctx server.MContext, err error) {
	log.Println("Error: Method("+ctx.Method()+"):", err)
	obj := bddp.NewError(ctx.Segment())
	obj.SetError(err.Error())
	ctx.SendError(&obj)
}
