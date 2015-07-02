package kmdb

import (
	"os"
	"reflect"
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

	o DatabaseConfig
)

// Test Server
// Start a test server with a test database and use it for
// tests performed later. Same database is used for all tests.

func init() {
	os.RemoveAll(DataPath)

	dc := DatabaseConfig{
		DataPath:       DataPath,
		IndexDepth:     4,
		PayloadSize:    16,
		BucketDuration: 3600000000000,
		Resolution:     60000000000,
		SegmentSize:    100000,
	}

	cfg := &ServerConfig{
		VerboseLogs:   true,
		RemoteDebug:   true,
		ListenAddress: Address,
		Databases: map[string]DatabaseConfig{
			Database: dc,
		},
	}

	dbs := map[string]kdb.Database{}
	db, err := dbase.New(dbase.Options{
		DatabaseName:   Database,
		DataPath:       dc.DataPath,
		IndexDepth:     dc.IndexDepth,
		PayloadSize:    dc.PayloadSize,
		BucketDuration: dc.BucketDuration,
		Resolution:     dc.Resolution,
		SegmentSize:    dc.SegmentSize,
	})

	if err != nil {
		panic(err)
	}

	dbs["test"] = db
	d = db
	o = dc

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

func TestEverything(t *testing.T) {
	ts := time.Now().UnixNano()
	f1 := []string{"a", "b", "c", "d"}
	f2 := []string{"a", "b", "c", "D"}
	g0 := []bool{true, true, true, true}

	// PUT DATA
	// --------

	b1 := &PutReqBatch{}
	b1.Batch = make([]*PutReq, 2, 2)

	b1.Batch[0] = &PutReq{
		Database:  Database,
		Timestamp: ts - o.Resolution*2,
		Value:     100.0,
		Count:     10,
		Fields:    f1,
	}

	b1.Batch[1] = &PutReq{
		Database:  Database,
		Timestamp: ts - o.Resolution*1,
		Value:     200.0,
		Count:     20,
		Fields:    f2,
	}

	o1, err := c.Put(b1)
	if err != nil {
		t.Fatal(err)
	}

	// TODO
	_ = o1

	// GET DATA (get)
	// --------------
	// verify put reqs

	b2 := &GetReqBatch{}
	b2.Batch = make([]*GetReq, 2, 2)

	b2.Batch[0] = &GetReq{
		Database:  Database,
		StartTime: ts - o.Resolution*2,
		EndTime:   ts - o.Resolution*1,
		Fields:    f1,
		GroupBy:   g0,
	}

	b2.Batch[1] = &GetReq{
		Database:  Database,
		StartTime: ts - o.Resolution*1,
		EndTime:   ts,
		Fields:    f2,
		GroupBy:   g0,
	}

	o2, err := c.Get(b2)
	if err != nil {
		t.Fatal(err)
	}

	if len(o2.Batch) != 2 {
		t.Fatal("incorrect result count")
	}

	// verify result no. 1
	o2r1 := o2.Batch[0]

	if len(o2r1.Data) != 1 {
		t.Fatal("incorrect series count")
	}

	o2r1s1 := o2r1.Data[0]

	if !reflect.DeepEqual(o2r1s1.Fields, f1) {
		t.Fatal("incorrect fields")
	}

	if len(o2r1s1.Points) != 1 {
		t.Fatal("incorrect number of points")
	}

	o2r1s1p1 := o2r1s1.Points[0]
	if o2r1s1p1.Count != 10 || !feq(o2r1s1p1.Value, 100.0) {
		t.Fatal("incorrect values")
	}

	// verify result no. 2
	o2r2 := o2.Batch[1]

	if len(o2r2.Data) != 1 {
		t.Fatal("incorrect series count")
	}

	o2r2s1 := o2r2.Data[0]

	if !reflect.DeepEqual(o2r2s1.Fields, f2) {
		t.Fatal("incorrect fields")
	}

	if len(o2r2s1.Points) != 1 {
		t.Fatal("incorrect number of points")
	}

	o2r2s1p1 := o2r2s1.Points[0]
	if o2r2s1p1.Count != 20 || !feq(o2r2s1p1.Value, 200.0) {
		t.Fatal("incorrect values")
	}

	// Inc DATA
	// --------

	b3 := &IncReqBatch{}
	b3.Batch = make([]*IncReq, 2, 2)

	b3.Batch[0] = &IncReq{
		Database:  Database,
		Timestamp: ts - o.Resolution*2,
		Value:     10.0,
		Count:     1,
		Fields:    f1,
	}

	b3.Batch[1] = &IncReq{
		Database:  Database,
		Timestamp: ts - o.Resolution*1,
		Value:     20.0,
		Count:     2,
		Fields:    f2,
	}

	o3, err := c.Inc(b3)
	if err != nil {
		t.Fatal(err)
	}

	// TODO
	_ = o3

	// GET DATA (get)
	// --------------
	// verify inc reqs

	// reuse batch b2
	b4 := b2
	o4, err := c.Get(b4)
	if err != nil {
		t.Fatal(err)
	}

	if len(o4.Batch) != 2 {
		t.Fatal("incorrect result count")
	}

	// verify result no. 1
	o4r1 := o4.Batch[0]

	if len(o4r1.Data) != 1 {
		t.Fatal("incorrect series count")
	}

	o4r1s1 := o4r1.Data[0]

	if !reflect.DeepEqual(o4r1s1.Fields, f1) {
		t.Fatal("incorrect fields")
	}

	if len(o4r1s1.Points) != 1 {
		t.Fatal("incorrect number of points")
	}

	o4r1s1p1 := o4r1s1.Points[0]
	if o4r1s1p1.Count != 11 || !feq(o4r1s1p1.Value, 110.0) {
		t.Fatal("incorrect values")
	}

	// verify result no. 2
	o4r2 := o4.Batch[1]

	if len(o4r2.Data) != 1 {
		t.Fatal("incorrect series count")
	}

	o4r2s1 := o4r2.Data[0]

	if !reflect.DeepEqual(o4r2s1.Fields, f2) {
		t.Fatal("incorrect fields")
	}

	if len(o4r2s1.Points) != 1 {
		t.Fatal("incorrect number of points")
	}

	o4r2s1p1 := o4r2s1.Points[0]
	if o4r2s1p1.Count != 22 || !feq(o4r2s1p1.Value, 220.0) {
		t.Fatal("incorrect values", o4r2s1p1)
	}
}

func feq(f1, f2 float64) bool {
	d := f1 - f2
	return d < 1.0 && d > -1.0
}
