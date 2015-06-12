package kmdb

import (
	"reflect"
	"testing"
	"time"

	"github.com/meteorhacks/kdb"
)

//   MockDB
// ----------

type MockDB struct {
	put_ts   int64
	put_vals []string
	put_pld  []byte

	get_start int64
	get_end   int64
	get_vals  []string

	find_start int64
	find_end   int64
	find_vals  []string
}

func (db *MockDB) Put(ts int64, vals []string, pld []byte) (err error) {
	db.put_ts = ts
	db.put_vals = vals
	db.put_pld = pld
	return nil
}

func (db *MockDB) Get(start, end int64, vals []string) (res [][]byte, err error) {
	db.get_start = start
	db.get_end = end
	db.get_vals = vals
	return nil, nil
}

func (db *MockDB) Find(start, end int64, vals []string) (res map[*kdb.IndexElement][][]byte, err error) {
	db.find_start = start
	db.find_end = end
	db.find_vals = vals
	return nil, nil
}

func (db *MockDB) RemoveBefore(ts int64) (err error) {
	return nil
}

func (db *MockDB) Close() (err error) {
	return nil
}

func (db *MockDB) Reset() {
	db.put_ts = 0
	db.put_vals = nil
	db.put_pld = nil
	db.get_start = 0
	db.get_end = 0
	db.get_vals = nil
	db.find_start = 0
	db.find_end = 0
	db.find_vals = nil
}

//   Init
// --------

var (
	dbs map[string]kdb.Database
	cfg *ServerConfig
	s   Server
	c   Client
)

// Test Server
// Start a test server with a mock database and use it for
// tests performed later. Mock DB is reset before each test.

func init() {
	cfg = &ServerConfig{
		DatabaseName:   "test",
		DataPath:       "/tmp/test",
		IndexDepth:     4,
		PayloadSize:    4,
		BucketDuration: 3600000000000,
		Resolution:     60000000000,
		SegmentSize:    100000,
		DebugMode:      true,
		BDDPAddress:    ":3000",
	}

	dbs = map[string]kdb.Database{
		"test": &MockDB{},
		"mock": &MockDB{},
	}

	s = NewServer(dbs, cfg)
	go s.Listen()

	// wait for the server to start
	time.Sleep(time.Second * 2)

	c = NewClient("localhost:3000")
	err := c.Connect()
	if err != nil {
		panic(err)
	}
}

func Reset() {
	for _, db := range dbs {
		db.(*MockDB).Reset()
	}
}

//   Tests
// ---------

func TestPut(t *testing.T) {
	defer Reset()

	b, err := c.PutBatch(1)
	if err != nil {
		t.Fatal(err)
	}

	var ts int64 = 123
	vals := []string{"a", "b", "c", "d"}
	pld := []byte{1, 2, 3, 4}

	err = b.Set("test", 0, ts, vals, pld)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Send()
	if err != nil {
		t.Fatal(err)
	}

	db := dbs["test"].(*MockDB)
	if db.put_ts != ts ||
		!reflect.DeepEqual(db.put_vals, vals) ||
		!reflect.DeepEqual(db.put_pld, pld) {
		t.Fatal("invalid value")
	}
}

func TestGet(t *testing.T) {
	defer Reset()

	b, err := c.GetBatch(1)
	if err != nil {
		t.Fatal(err)
	}

	var start int64 = 10
	var end int64 = 20
	vals := []string{"a", "b", "c", "d"}

	err = b.Set("test", 0, vals, start, end)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Send()
	if err != nil {
		t.Fatal(err)
	}

	db := dbs["test"].(*MockDB)
	if db.get_start != start ||
		db.get_end != end ||
		!reflect.DeepEqual(db.get_vals, vals) {
		t.Fatal("invalid value")
	}
}

func TestFind(t *testing.T) {
	defer Reset()

	b, err := c.GetBatch(1)
	if err != nil {
		t.Fatal(err)
	}

	var start int64 = 10
	var end int64 = 20
	vals := []string{"a", "b", "c", ""}

	err = b.Set("test", 0, vals, start, end)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Send()
	if err != nil {
		t.Fatal(err)
	}

	db := dbs["test"].(*MockDB)
	if db.find_start != start ||
		db.find_end != end ||
		!reflect.DeepEqual(db.find_vals, vals) {
		t.Fatal("invalid value")
	}
}
