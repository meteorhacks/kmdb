package kmdb

import (
	"fmt"

	"github.com/meteorhacks/kdb"
)

type MockDB struct {
}

func (db *MockDB) Put(ts int64, vals []string, pld []byte) (err error) {
	fmt.Println("Put:", ts, vals, pld)

}

func (db *MockDB) Get(start, end int64, vals []string) (res [][]byte, err error) {
	fmt.Println("Get:", start, end, vals)

}

func (db *MockDB) Find(start, end int64, vals []string) (res map[*kdb.IndexElement][][]byte, err error) {
	fmt.Println("Find:", start, end, vals)

}

func (db *MockDB) RemoveBefore(ts int64) (err error) {
	fmt.Println("RemoveBefore:", ts)

}

func (db *MockDB) Close() (err error) {
	fmt.Println("Close:")

}
