package kmdb

import (
	"log"

	"github.com/glycerine/go-capnproto"
	"github.com/meteorhacks/bddp"
	"github.com/meteorhacks/bddp/server"
	"github.com/meteorhacks/kdb"
)

type Server struct {
	Database kdb.Database
	server   server.Server
	config   *Config
}

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

	params := PutRequest_List(*ctx.Params())
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
	res := NewPutResult(seg)
	res.SetOk(true)

	obj := capn.Object(res)
	ctx.SendResult(&obj)
}

// Method = "get"
// receives a `GetRequest` list and responds with a list of `GetResult`
// uses either `db.Get()` or `db.Find()` to get data from the database
func (s *Server) handleGet(ctx server.MContext) {
	defer ctx.SendUpdated()

	params := GetRequest_List(*ctx.Params())
	count := params.Len()

	seg := ctx.Segment()
	out := NewGetResultList(seg, count)

	for i := 0; i < count; i++ {
		req := params.At(i)
		start := req.Start()
		end := req.End()
		vals := req.Values().ToArray()

		res := NewGetResult(seg)
		out.Set(i, res)

		var resItems ResultItem_List

		// use the `Get` method only if all values are set
		// otherwise use the more costly `Find` method
		if s.canUseGet(vals) {
			data, err := s.Database.Get(start, end, vals)
			if err != nil {
				log.Println(err)
				continue
			}

			resItems = NewResultItemList(seg, 1)
			resItem := s.makeResultItem(seg, data, vals)
			resItems.Set(0, *resItem)
		} else {
			dataMap, err := s.Database.Find(start, end, vals)
			if err != nil {
				log.Println(err)
				continue
			}

			numItems := len(dataMap)
			resItems = NewResultItemList(seg, numItems)

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

func (s *Server) makeResultItem(seg *capn.Segment, data [][]byte, vals []string) (resItem *ResultItem) {
	item := NewResultItem(seg)
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
