package kmdb

import (
	"errors"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/meteorhacks/kdb"
	"golang.org/x/net/context"
)

const (
	PayloadSize = 16
)

var (
	ErrBatchError = errors.New("batch didn't complete successfully")
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
	ListenAddress string `json:"listen_address"`

	Databases map[string]DatabaseConfig `json:"databases"`
}

//   Server
// ----------

type Server interface {
	DatabaseServiceServer
	Listen() (err error)
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
	lis, err := net.Listen("tcp", s.cfg.ListenAddress)
	if err != nil {
		return err
	}

	gsrv := grpc.NewServer()
	RegisterDatabaseServiceServer(gsrv, s)

	log.Println("GRPCS:  listening on", s.cfg.ListenAddress)
	return gsrv.Serve(lis)
}

func (s *server) Put(ctx context.Context, req *PutReq) (r *PutRes, err error) {
	return s.put(req)
}

func (s *server) Inc(ctx context.Context, req *IncReq) (r *IncRes, err error) {
	return s.inc(req)
}

func (s *server) Get(ctx context.Context, req *GetReq) (r *GetRes, err error) {
	return s.get(req)
}

func (s *server) PutBatch(ctx context.Context, batch *PutReqBatch) (r *PutResBatch, berr error) {
	n := len(batch.GetBatch())
	r = &PutResBatch{}
	r.Batch = make([]*PutRes, n, n)
	var err error

	for i := 0; i < n; i++ {
		r.Batch[i], err = s.put(batch.Batch[i])
		if err != nil && berr == nil {
			berr = ErrBatchError
		}
	}

	return r, berr
}

func (s *server) IncBatch(ctx context.Context, batch *IncReqBatch) (r *IncResBatch, berr error) {
	n := len(batch.GetBatch())
	r = &IncResBatch{}
	r.Batch = make([]*IncRes, n, n)
	var err error

	for i := 0; i < n; i++ {
		r.Batch[i], err = s.inc(batch.Batch[i])
		if err != nil && berr == nil {
			berr = ErrBatchError
		}
	}

	return r, berr
}

func (s *server) GetBatch(ctx context.Context, batch *GetReqBatch) (r *GetResBatch, berr error) {
	n := len(batch.GetBatch())
	r = &GetResBatch{}
	r.Batch = make([]*GetRes, n, n)
	var err error

	for i := 0; i < n; i++ {
		r.Batch[i], err = s.get(batch.Batch[i])
		if err != nil && berr == nil {
			berr = ErrBatchError
		}
	}

	return r, berr
}

func (s *server) put(req *PutReq) (r *PutRes, err error) {
	r = &PutRes{}
	return r, nil
}

func (s *server) inc(req *IncReq) (r *IncRes, err error) {
	r = &IncRes{}
	return r, nil
}

func (s *server) get(req *GetReq) (r *GetRes, err error) {
	r = &GetRes{}
	return r, nil
}
