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

var (
	d kdb.Database
	s Server
	c Client
)

// Test Server
// Start a test server with a test database and use it for
// tests performed later. Same database is used for all tests.

func init() {
	os.RemoveAll("/tmp/kmdb_test")

	dcfg := DatabaseConfig{
		DataPath:       "/tmp/kmdb_test",
		IndexDepth:     4,
		PayloadSize:    16,
		BucketDuration: 3600000000000,
		Resolution:     60000000000,
		SegmentSize:    100000,
	}

	cfg := &ServerConfig{
		RemoteDebug:   true,
		ListenAddress: "localhost:3000",
		Databases: map[string]DatabaseConfig{
			"test": dcfg,
		},
	}

	dbs := map[string]kdb.Database{}
	db, err := dbase.New(dbase.Options{
		DatabaseName:   "test",
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

	c = NewClient("localhost:3000")
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
