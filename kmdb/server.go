package kmdb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/meteorhacks/kdb"
	"github.com/meteorhacks/simple-rpc-go"
)

var (
	ErrDBNotFound = errors.New("requested db is not setup on this server")
	ErrBatchError = errors.New("batch didn't complete successfully")
)

type DatabaseConfig struct {
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

	// address to listen for thrift traffic (host:port)
	ListenAddress string `json:"listen_address"`

	Databases map[string]DatabaseConfig `json:"databases"`
}

//   Server
// ----------

type Server interface {
	Listen() (err error)
	Put(req []byte) (res []byte, err error)
	Inc(req []byte) (res []byte, err error)
	Get(req []byte) (res []byte, err error)
}

type server struct {
	cfg *ServerConfig
	dbs map[string]kdb.Database
}

func NewServer(dbs map[string]kdb.Database, cfg *ServerConfig) (s Server) {
	ss := &server{cfg, dbs}
	return ss
}

func (s *server) Listen() (err error) {
	srv := srpc.NewServer(s.cfg.ListenAddress)
	srv.SetHandler("put", s.Put)
	srv.SetHandler("inc", s.Inc)
	srv.SetHandler("get", s.Get)

	log.Println("SRPCS:  listening on", s.cfg.ListenAddress)
	return srv.Listen()
}

func (s *server) Put(req []byte) (res []byte, err error) {
	batch := &PutReqBatch{}
	err = proto.Unmarshal(req, batch)
	if err != nil {
		return nil, err
	}

	n := len(batch.Batch)
	r := &PutResBatch{}
	r.Batch = make([]*PutRes, n, n)
	var batchError error

	for i := 0; i < n; i++ {
		r.Batch[i], err = s.put(batch.Batch[i])
		if err != nil && batchError == nil {
			batchError = ErrBatchError
		}
	}

	if batchError != nil {
		return nil, batchError
	}

	res, err = proto.Marshal(r)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *server) Inc(req []byte) (res []byte, err error) {
	batch := &IncReqBatch{}
	err = proto.Unmarshal(req, batch)
	if err != nil {
		return nil, err
	}

	n := len(batch.Batch)
	r := &IncResBatch{}
	r.Batch = make([]*IncRes, n, n)
	var batchError error

	for i := 0; i < n; i++ {
		r.Batch[i], err = s.inc(batch.Batch[i])
		if err != nil && batchError == nil {
			batchError = ErrBatchError
		}
	}

	if batchError != nil {
		return nil, batchError
	}

	res, err = proto.Marshal(r)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *server) Get(req []byte) (res []byte, err error) {
	batch := &GetReqBatch{}
	err = proto.Unmarshal(req, batch)
	if err != nil {
		return nil, err
	}

	n := len(batch.Batch)
	r := &GetResBatch{}
	r.Batch = make([]*GetRes, n, n)
	var batchError error

	for i := 0; i < n; i++ {
		r.Batch[i], err = s.get(batch.Batch[i])
		if err != nil && batchError == nil {
			batchError = ErrBatchError
		}
	}

	if batchError != nil {
		return nil, batchError
	}

	res, err = proto.Marshal(r)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *server) put(req *PutReq) (r *PutRes, err error) {
	r = &PutRes{}

	db, _, err := s.getDB(req.Database)
	if err != nil {
		return r, err
	}

	pld := valToPld(req.Value, req.Count)
	if err := db.Put(req.Timestamp, req.Fields, pld); err != nil {
		return r, err
	}

	return r, nil
}

func (s *server) inc(req *IncReq) (r *IncRes, err error) {
	r = &IncRes{}

	db, dbCfg, err := s.getDB(req.Database)
	if err != nil {
		return r, err
	}

	ts1 := req.Timestamp
	ts2 := ts1 + dbCfg.Resolution
	out, err := db.Get(ts1, ts2, req.Fields)
	if err != nil {
		return r, err
	}

	val, num := pldToVal(out[0])
	val += req.Value
	num += req.Count
	pld := valToPld(val, num)

	if err := db.Put(ts1, req.Fields, pld); err != nil {
		return r, err
	}

	return r, nil
}

func (s *server) get(req *GetReq) (r *GetRes, err error) {
	r = &GetRes{}

	db, dbCfg, err := s.getDB(req.Database)
	if err != nil {
		return r, err
	}

	ts1 := req.StartTime
	ts2 := req.EndTime
	fields := req.Fields
	groupBy := req.GroupBy

	gettable := true
	vcount := int(dbCfg.IndexDepth)
	for j := 0; j < vcount; j++ {
		if fields[j] == "" {
			gettable = false
		}
	}

	var ss *seriesSet
	// use the `Get` method only if all values are set
	// otherwise use the more costly `Find` method
	if gettable {
		ss, err = s.getWithGet(db, ts1, ts2, fields, groupBy)
	} else {
		ss, err = s.getWithFind(db, ts1, ts2, fields, groupBy)
	}

	if err != nil {
		return r, err
	}

	r.Data = ss.toResult()

	return r, nil
}

// Get database and database config
func (s *server) getDB(name string) (db kdb.Database, cfg *DatabaseConfig, err error) {
	db, ok := s.dbs[name]
	if !ok {
		return nil, nil, ErrDBNotFound
	}

	config, ok := s.cfg.Databases[name]
	if !ok {
		return nil, nil, ErrDBNotFound
	}

	return db, &config, nil
}

func (s *server) getWithGet(db kdb.Database, start, end int64, fields []string, groupBy []bool) (ss *seriesSet, err error) {
	data, err := db.Get(start, end, fields)
	if err != nil {
		return nil, err
	}

	ss = s.newSeriesSet(groupBy)
	sr := s.newSeries(data, fields)
	ss.add(sr)

	return ss, nil
}

func (s *server) getWithFind(db kdb.Database, start, end int64, fields []string, groupBy []bool) (ss *seriesSet, err error) {
	dataMap, err := db.Find(start, end, fields)
	if err != nil {
		return nil, err
	}

	ss = s.newSeriesSet(groupBy)

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

func (s *server) newSeriesSet(groupBy []bool) (ss *seriesSet) {
	set := []*series{}
	return &seriesSet{set, groupBy}
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

func (p *point) toResult() (item *ResPoint) {
	item = &ResPoint{}
	item.Value = p.value
	item.Count = p.count
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

func (sr *series) toResult() (item *ResSeries) {
	item = &ResSeries{}
	item.Fields = sr.fields

	count := len(sr.points)
	item.Points = make([]*ResPoint, count, count)
	for i, p := range sr.points {
		point := p.toResult()
		item.Points[i] = point
	}

	return item
}

type seriesSet struct {
	items   []*series
	groupBy []bool
}

func (ss *seriesSet) add(sn *series) {
	ss.grpFields(sn)

	count := len(ss.items)
	for i := 0; i < count; i++ {
		sr := ss.items[i]
		if sr.canMerge(sn) {
			sr.add(sn)
			return
		}
	}

	ss.items = append(ss.items, sn)
}

func (ss *seriesSet) grpFields(sn *series) {
	count := len(sn.fields)
	grouped := make([]string, count, count)

	for i := 0; i < count; i++ {
		if ss.groupBy[i] {
			grouped[i] = sn.fields[i]
		}
	}

	sn.fields = grouped
}

func (ss *seriesSet) toResult() (res []*ResSeries) {
	count := len(ss.items)
	res = make([]*ResSeries, count, count)

	for i := 0; i < count; i++ {
		sr := ss.items[i]
		item := sr.toResult()
		res[i] = item
	}

	return res
}
