package kmdb

import (
	"log"

	"github.com/glycerine/go-capnproto"
	"github.com/meteorhacks/bddp"
	"github.com/meteorhacks/kdb"
)

type ServerConfig struct {
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

//   Server
// ----------

type Server interface {
	Listen() (err error)
}

type server struct {
	*ServerConfig
	db kdb.Database
	bs bddp.Server
}

func NewServer(db kdb.Database, cfg *ServerConfig) (s Server) {
	bs := bddp.NewServer(cfg.BDDPAddress)
	ss := &server{cfg, db, bs}

	// method handlers
	bs.Method("put", ss.handlePut)
	bs.Method("get", ss.handleGet)

	return ss
}

func (s *server) Listen() (err error) {
	return s.bs.Listen()
}

// Method = "put"
// receives a `PutRequest` list and saves metrics 1 by 1
// on success sends a `PutResult` with `ok` set to true
func (s *server) handlePut(ctx bddp.MContext) {
	defer ctx.SendUpdated()

	params := PutRequest_List(*ctx.Params())
	pcount := params.Len()

	vcount := int(s.IndexDepth)
	vals := make([]string, vcount, vcount)
	seg := ctx.Segment()

	for i := 0; i < pcount; i++ {
		req := params.At(i)
		ts := req.Time()
		pld := req.Payload()

		vals := vals[:0]
		for j := 0; j < vcount; j++ {
			vals = append(vals, req.Values().At(j))
		}

		if err := s.db.Put(ts, vals, pld); err != nil {
			s.handleErr(ctx, err)
			return
		}
	}

	res := NewPutResult(seg)
	res.SetOk(true)

	obj := capn.Object(res)
	ctx.SendResult(&obj)
}

// Method = "get"
// receives a `GetRequest` list and responds with a list of `GetResult`
// uses either `db.Get()` or `db.Find()` to get data from the database
func (s *server) handleGet(ctx bddp.MContext) {
	defer ctx.SendUpdated()

	params := GetRequest_List(*ctx.Params())
	pcount := params.Len()

	vcount := int(s.IndexDepth)
	vals := make([]string, vcount, vcount)
	seg := ctx.Segment()
	ress := NewGetResultList(seg, pcount)

	for i := 0; i < pcount; i++ {
		req := params.At(i)
		start := req.Start()
		end := req.End()

		vals := vals[:0]
		gettable := true
		for j := 0; j < vcount; j++ {
			v := req.Values().At(j)
			vals = append(vals, v)

			if v == "" {
				gettable = false
			}
		}

		var items ResultItem_List

		// use the `Get` method only if all values are set
		// otherwise use the more costly `Find` method
		if gettable {
			data, err := s.db.Get(start, end, vals)
			if err != nil {
				s.handleErr(ctx, err)
				return
			}

			items = NewResultItemList(seg, 1)
			item := newResultItem(seg, data, vals)
			items.Set(0, *item)
		} else {
			dataMap, err := s.db.Find(start, end, vals)
			if err != nil {
				s.handleErr(ctx, err)
				return
			}

			numItems := len(dataMap)
			items = NewResultItemList(seg, numItems)

			counter := 0
			for el, data := range dataMap {
				item := newResultItem(seg, data, el.Values)
				items.Set(counter, *item)
				counter++
			}
		}

		res := NewGetResult(seg)
		ress.Set(i, res)
		res.SetOk(true)
		res.SetData(items)
	}

	obj := capn.Object(ress)
	ctx.SendResult(&obj)
}

// Sends a method error
// Converts a go error to a method error and send.
// `handleErr` can be used by any method handlers.
func (s *server) handleErr(ctx bddp.MContext, err error) {
	log.Println("Error: Method("+ctx.Method()+"):", err)
	obj := bddp.NewError(ctx.Segment())
	obj.SetError(err.Error())
	ctx.SendError(&obj)
}

// Creates a ResultItem
func newResultItem(seg *capn.Segment, pld [][]byte, vals []string) (item *ResultItem) {
	itm := NewResultItem(seg)
	item = &itm

	dlist := seg.NewDataList(len(pld))
	item.SetData(dlist)
	for j, d := range pld {
		dlist.Set(j, d)
	}

	vlist := seg.NewTextList(len(vals))
	item.SetValues(vlist)
	for j, v := range vals {
		vlist.Set(j, v)
	}

	return item
}
