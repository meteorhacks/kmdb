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
	resolution int64

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
	count := int(end-start) / 10
	res = makePayload(count, 1)
	return res, nil
}

func (db *MockDB) Find(start, end int64, vals []string) (res map[*kdb.IndexElement][][]byte, err error) {
	db.find_start = start
	db.find_end = end
	db.find_vals = vals

	count := int(end-start) / 10
	res = make(map[*kdb.IndexElement][][]byte)

	el1 := &kdb.IndexElement{Values: []string{"a", "b", "c", "d"}}
	res[el1] = makePayload(count, 2)

	el2 := &kdb.IndexElement{Values: []string{"a", "b", "c", "e"}}
	res[el2] = makePayload(count, 3)

	return res, nil
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

func makePayload(count, mult int) (res [][]byte) {
	res = make([][]byte, count, count)
	for i := 0; i < count; i++ {
		res[i] = valToPld(float64(mult*(i+1)*10), int64(mult*(i+1)))
	}

	return res
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
		RemoteDebug: true,
		BDDPAddress: ":3000",
		Databases: map[string]DatabaseConfig{
			"test": DatabaseConfig{
				DatabaseName:   "test",
				DataPath:       "/tmp/test",
				IndexDepth:     4,
				PayloadSize:    16,
				BucketDuration: 1000,
				Resolution:     10,
				SegmentSize:    100,
			},
			"mock": DatabaseConfig{
				DatabaseName:   "mock",
				DataPath:       "/tmp/mock",
				IndexDepth:     4,
				PayloadSize:    16,
				BucketDuration: 1000,
				Resolution:     10,
				SegmentSize:    100,
			},
		},
	}

	dbs = map[string]kdb.Database{
		"test": &MockDB{resolution: 10},
		"mock": &MockDB{resolution: 10},
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

	val := 123.45
	var num int64 = 67890
	pld := valToPld(val, num)

	err = b.Set(0, "test", ts, vals, val, num)
	if err != nil {
		t.Fatal(err)
	}

	_, err = b.Send()
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

func TestInc(t *testing.T) {
	defer Reset()

	var ts int64 = 123
	vals := []string{"a", "b", "c", "d"}

	val := 123.0
	var num int64 = 10

	b, err := c.IncBatch(1)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Set(0, "test", ts, vals, val, num)
	if err != nil {
		t.Fatal(err)
	}

	_, err = b.Send()
	if err != nil {
		t.Fatal(err)
	}

	// MockDB gives out (10,1) as first value
	// it should be added with the new value
	pld := valToPld(10.0+123.0, 10+1)

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

	var start int64 = 100
	var end int64 = 200
	vals := []string{"a", "b", "c", "d"}
	grup := []bool{true, true, false, true}

	err = b.Set(0, "test", start, end, vals, grup)
	if err != nil {
		t.Fatal(err)
	}

	obj, err := b.Send()
	if err != nil {
		t.Fatal(err)
	}

	ress := GetResult_List(obj)
	if ress.Len() != 1 {
		t.Fatal(err)
	}

	res := ress.At(0)

	ss := res.Data()
	if ss.Len() != 1 {
		t.Fatal(err)
	}

	sr := ss.At(0)
	fields := sr.Fields().ToArray()
	expFields := []string{"a", "b", "", "d"}
	if !reflect.DeepEqual(fields, expFields) {
		t.Fatal("incorrect fields")
	}

	points := sr.Points()
	expPoints := makePayload(10, 1)
	if points.Len() != 10 {
		t.Fatal("incorrect points count")
	}

	for i := 0; i < 10; i++ {
		rp := points.At(i)
		val, num := pldToVal(expPoints[i])
		if val != rp.Value() || num != rp.Count() {
			t.Fatal("incorrect value")
		}
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

	var start int64 = 100
	var end int64 = 200
	vals := []string{"a", "b", "c", ""}
	grup := []bool{true, true, true, false}

	err = b.Set(0, "test", start, end, vals, grup)
	if err != nil {
		t.Fatal(err)
	}

	obj, err := b.Send()
	if err != nil {
		t.Fatal(err)
	}

	ress := GetResult_List(obj)
	if ress.Len() != 1 {
		t.Fatal(err)
	}

	res := ress.At(0)

	ss := res.Data()
	if ss.Len() != 1 {
		t.Fatal(err)
	}

	sr := ss.At(0)
	fields := sr.Fields().ToArray()
	expFields := []string{"a", "b", "c", ""}
	if !reflect.DeepEqual(fields, expFields) {
		t.Fatal("incorrect fields")
	}

	points := sr.Points()
	expPoints := makePayload(10, 5)
	if points.Len() != 10 {
		t.Fatal("incorrect points count")
	}

	for i := 0; i < 10; i++ {
		rp := points.At(i)
		val, num := pldToVal(expPoints[i])
		if val != rp.Value() || num != rp.Count() {
			t.Fatal("incorrect value")
		}
	}

	db := dbs["test"].(*MockDB)
	if db.find_start != start ||
		db.find_end != end ||
		!reflect.DeepEqual(db.find_vals, vals) {
		t.Fatal("invalid value")
	}
}
