package kmdb

import (
	"os"
	"testing"
	"time"

	"github.com/meteorhacks/kdb"
	"github.com/meteorhacks/kdb/dbase"
)

//   Init
// --------

const (
	Address  = "localhost:3000"
	DataPath = "/tmp/kmdb_test"
	Database = "test"
)

var (
	d kdb.Database
	s Server
	c Client
)

// Test Server
// Start a test server with a test database and use it for
// tests performed later. Same database is used for all tests.

func init() {
	os.RemoveAll(DataPath)

	dcfg := DatabaseConfig{
		DataPath:       DataPath,
		IndexDepth:     4,
		PayloadSize:    16,
		BucketDuration: 3600000000000,
		Resolution:     60000000000,
		SegmentSize:    100000,
	}

	cfg := &ServerConfig{
		RemoteDebug:   true,
		ListenAddress: Address,
		Databases: map[string]DatabaseConfig{
			Database: dcfg,
		},
	}

	dbs := map[string]kdb.Database{}
	db, err := dbase.New(dbase.Options{
		DatabaseName:   Database,
		DataPath:       dcfg.DataPath,
		IndexDepth:     dcfg.IndexDepth,
		PayloadSize:    dcfg.PayloadSize,
		BucketDuration: dcfg.BucketDuration,
		Resolution:     dcfg.Resolution,
		SegmentSize:    dcfg.SegmentSize,
	})

	if err != nil {
		panic(err)
	}

	dbs["test"] = db
	d = db

	s = NewServer(dbs, cfg)
	go s.Listen()

	// wait for the server to start
	time.Sleep(time.Second * 2)

	c = NewClient(Address)
	if err := c.Connect(); err != nil {
		panic(err)
	}
}

//   Tests
// ---------

func TestPut(t *testing.T) {
	// TODO: write test
}

func TestInc(t *testing.T) {
	// TODO: write test
}

func TestGet(t *testing.T) {
	// TODO: write test
}
