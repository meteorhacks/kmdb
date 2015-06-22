package kmdb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"

	"github.com/glycerine/go-capnproto"
	"github.com/meteorhacks/bddp"
	"github.com/meteorhacks/kdb"
)

const (
	PayloadSize = 16
)

var (
	ErrDBNotFound = errors.New("db not found")
)

type DatabaseConfig struct {
	// database name. Currently only used with naming files
	// can be useful when supporting multiple Databases
	DatabaseName string `json:"database_name"`

	// place to store data files
	DataPath string `json:"database_path"`

	// depth of the index tree
	IndexDepth int64 `json:"index_depth"`

	// payload size should always be equal to this amount
	PayloadSize int64 `json:"payload_size"`

	// bucket resolution in nano seconds
	Resolution int64 `json:"payload_resolution"`

	// time duration in nano seconds of a range unit
	// this should be a multiple of `Resolution`
	BucketDuration int64 `json:"bucket_duration"`

	// number of records per segment
	SegmentSize int64 `json:"segment_size"`
}

type ServerConfig struct {
	// enable pprof on ":6060" instead of "localhost:6060".
	RemoteDebug bool `json:"remote_debug"`

	// address to listen for bddp traffic (host:port)
	BDDPAddress string `json:"bddp_address"`

	Databases map[string]DatabaseConfig `json:"databases"`
}

//   Server
// ----------

type Server interface {
	Listen() (err error)
}

type server struct {
	*ServerConfig
	dbs map[string]kdb.Database
	bs  bddp.Server
}

func NewServer(dbs map[string]kdb.Database, cfg *ServerConfig) (s Server) {
	bs := bddp.NewServer(cfg.BDDPAddress)
	ss := &server{cfg, dbs, bs}

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

	vals := []string{}
	seg := ctx.Segment()

	for i := 0; i < pcount; i++ {
		req := params.At(i)

		dbName := req.Database()
		db, dbCfg, err := s.getDB(dbName)
		if err != nil {
			s.handleErr(ctx, err)
			return
		}

		ts := req.Timestamp()
		pld := valToPld(req.Value(), req.Count())

		vals := vals[:0]
		vcount := int(dbCfg.IndexDepth)
		for j := 0; j < vcount; j++ {
			vals = append(vals, req.Fields().At(j))
		}

		if err := db.Put(ts, vals, pld); err != nil {
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

	fields := []string{}
	groupBy := []bool{}
	seg := ctx.Segment()
	ress := NewGetResultList(seg, pcount)

	for i := 0; i < pcount; i++ {
		req := params.At(i)

		dbName := req.Database()
		db, dbCfg, err := s.getDB(dbName)
		if err != nil {
			s.handleErr(ctx, err)
			return
		}

		start := req.StartTime()
		end := req.EndTime()

		fields := fields[:0]
		groupBy := groupBy[:0]
		gettable := true

		vcount := int(dbCfg.IndexDepth)
		for j := 0; j < vcount; j++ {
			v := req.Fields().At(j)
			fields = append(fields, v)

			b := req.GroupBy().At(j)
			groupBy = append(groupBy, b)

			if v == "" {
				gettable = false
			}
		}

		var ss *seriess

		// use the `Get` method only if all values are set
		// otherwise use the more costly `Find` method
		if gettable {
			ss, err = s.getData(db, start, end, fields, groupBy)
		} else {
			ss, err = s.findData(db, start, end, fields, groupBy)
		}

		if err != nil {
			s.handleErr(ctx, err)
			return
		}

		items := ss.toResult(seg)
		res := NewGetResult(seg)
		res.SetOk(true)
		res.SetData(items)
		ress.Set(i, res)
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

// Get database and database config
func (s *server) getDB(name string) (db kdb.Database, cfg *DatabaseConfig, err error) {
	db, ok := s.dbs[name]
	if !ok {
		return nil, nil, ErrDBNotFound
	}

	config, ok := s.Databases[name]
	if !ok {
		return nil, nil, ErrDBNotFound
	}

	return db, &config, nil
}

func (s *server) getData(db kdb.Database, start, end int64, fields []string, groupBy []bool) (ss *seriess, err error) {
	data, err := db.Get(start, end, fields)
	if err != nil {
		return nil, err
	}

	ss = s.newSeriess(groupBy)

	sr := s.newSeries(data, fields)
	ss.add(sr)

	return ss, nil
}

func (s *server) findData(db kdb.Database, start, end int64, fields []string, groupBy []bool) (ss *seriess, err error) {
	dataMap, err := db.Find(start, end, fields)
	if err != nil {
		return nil, err
	}

	ss = s.newSeriess(groupBy)

	for el, data := range dataMap {
		sr := s.newSeries(data, el.Values)
		ss.add(sr)
	}

	return ss, nil
}

func (s *server) newSeries(data [][]byte, fields []string) (sr *series) {
	count := len(data)
	points := make([]*point, count, count)

	for i := 0; i < count; i++ {
		val, num := pldToVal(data[i])
		points[i] = &point{val, num}
	}

	return &series{fields, points, data}
}

func (s *server) newSeriess(groupBy []bool) (ss *seriess) {
	return &seriess{make([]*series, 0, 1), groupBy}
}

func valToPld(val float64, num int64) (pld []byte) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, val)
	binary.Write(buf, binary.LittleEndian, num)
	return buf.Bytes()
}

func pldToVal(pld []byte) (val float64, num int64) {
	buf := bytes.NewBuffer(pld)
	binary.Read(buf, binary.LittleEndian, &val)
	binary.Read(buf, binary.LittleEndian, &num)
	return val, num
}

// Helper structs for building get results

type point struct {
	value float64
	count int64
}

func (p *point) add(q *point) {
	p.value += q.value
	p.count += q.count
}

func (p *point) toResult(seg *capn.Segment) (item *ResultPoint) {
	itm := NewResultPoint(seg)
	item = &itm
	item.SetCount(p.count)
	item.SetValue(p.value)
	return item
}

type series struct {
	fields []string
	points []*point
	data   [][]byte
}

func (sr *series) add(sn *series) {
	count := len(sr.points)
	for i := 0; i < count; i++ {
		sr.points[i].add(sn.points[i])
	}
}

func (sr *series) canMerge(sn *series) (can bool) {
	count := len(sr.fields)
	for i := 0; i < count; i++ {
		if sr.fields[i] != sn.fields[i] {
			return false
		}
	}

	return true
}

func (sr *series) toResult(seg *capn.Segment) (item *ResultSeries) {
	itm := NewResultSeries(seg)
	item = &itm

	count := len(sr.points)
	points := NewResultPointList(seg, count)
	item.SetPoints(points)
	for j, p := range sr.points {
		point := p.toResult(seg)
		points.Set(j, *point)
	}

	fields := seg.NewTextList(len(sr.fields))
	item.SetFields(fields)
	for j, v := range sr.fields {
		fields.Set(j, v)
	}

	return item
}

type seriess struct {
	seriess []*series
	groupBy []bool
}

func (ss *seriess) add(sn *series) {
	ss.grpFields(sn)

	count := len(ss.seriess)
	for i := 0; i < count; i++ {
		sr := ss.seriess[i]
		if sr.canMerge(sn) {
			sr.add(sn)
			return
		}
	}

	ss.seriess = append(ss.seriess, sn)
}

func (ss *seriess) grpFields(sn *series) {
	count := len(sn.fields)
	grouped := make([]string, count, count)

	for i := 0; i < count; i++ {
		if ss.groupBy[i] {
			grouped[i] = sn.fields[i]
		}
	}

	sn.fields = grouped
}

func (ss *seriess) toResult(seg *capn.Segment) (res ResultSeries_List) {
	count := len(ss.seriess)
	res = NewResultSeriesList(seg, count)

	for i := 0; i < count; i++ {
		sr := ss.seriess[i]
		item := sr.toResult(seg)
		res.Set(i, *item)
	}

	return res
}
