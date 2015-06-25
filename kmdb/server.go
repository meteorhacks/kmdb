package kmdb

import (
	"log"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/meteorhacks/kdb"
)

const (
	PayloadSize = 16
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

	// address to listen for thrift traffic (host:port)
	ThritAddress string `json:"thrift_address"`

	Databases map[string]DatabaseConfig `json:"databases"`
}

//   Server
// ----------

type Server interface {
	ThriftService
	Listen() (err error)
}

type server struct {
	*ServerConfig
	dbs map[string]kdb.Database
}

func NewServer(dbs map[string]kdb.Database, cfg *ServerConfig) (s Server) {
	ss := &server{cfg, dbs}
	return ss
}

func (s *server) Listen() (err error) {
	tfac := thrift.NewTTransportFactory()
	// pfac := thrift.NewTBinaryProtocolFactoryDefault()
	// pfac := thrift.NewTCompactProtocolFactory()
	pfac := thrift.NewTJSONProtocolFactory()

	trans, err := thrift.NewTServerSocket(s.ThritAddress)
	if err != nil {
		return err
	}

	proc := NewThriftServiceProcessor(s)
	server := thrift.NewTSimpleServer4(proc, trans, tfac, pfac)

	log.Println("THRIFT: listening on", s.ThritAddress)
	return server.Serve()
}

func (s *server) Put(req *PutReq) (r *PutRes, err error) {
	return s.put(req)
}

func (s *server) Inc(req *IncReq) (r *IncRes, err error) {
	return s.inc(req)
}

func (s *server) Get(req *GetReq) (r *GetRes, err error) {
	return s.get(req)
}

func (s *server) PutBatch(batch []*PutReq) (r []*PutRes, berr error) {
	n := len(batch)
	r = make([]*PutRes, n, n)
	var err error

	for i := 0; i < n; i++ {
		r[i], err = s.put(batch[i])
		if err != nil && berr == nil {
			berr = ERR_BATCH_ERROR
		}
	}

	return r, berr
}

func (s *server) IncBatch(batch []*IncReq) (r []*IncRes, berr error) {
	n := len(batch)
	r = make([]*IncRes, n, n)
	var err error

	for i := 0; i < n; i++ {
		r[i], err = s.inc(batch[i])
		if err != nil && berr == nil {
			berr = ERR_BATCH_ERROR
		}
	}

	return r, berr
}

func (s *server) GetBatch(batch []*GetReq) (r []*GetRes, berr error) {
	n := len(batch)
	r = make([]*GetRes, n, n)
	var err error

	for i := 0; i < n; i++ {
		r[i], err = s.get(batch[i])
		if err != nil && berr == nil {
			berr = ERR_BATCH_ERROR
		}
	}

	return r, berr
}

func (s *server) put(req *PutReq) (r *PutRes, err error) {
	r = NewPutRes()
	return r, nil
}

func (s *server) inc(req *IncReq) (r *IncRes, err error) {
	r = NewIncRes()
	return r, nil
}

func (s *server) get(req *GetReq) (r *GetRes, err error) {
	r = NewGetRes()
	return r, nil
}
