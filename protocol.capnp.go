package kmdb

// AUTO GENERATED - DO NOT EDIT

import (
	C "github.com/glycerine/go-capnproto"
)

type PutRequest C.Struct

func NewPutRequest(s *C.Segment) PutRequest      { return PutRequest(s.NewStruct(16, 2)) }
func NewRootPutRequest(s *C.Segment) PutRequest  { return PutRequest(s.NewRootStruct(16, 2)) }
func AutoNewPutRequest(s *C.Segment) PutRequest  { return PutRequest(s.NewStructAR(16, 2)) }
func ReadRootPutRequest(s *C.Segment) PutRequest { return PutRequest(s.Root(0).ToStruct()) }
func (s PutRequest) Partition() int64            { return int64(C.Struct(s).Get64(0)) }
func (s PutRequest) SetPartition(v int64)        { C.Struct(s).Set64(0, uint64(v)) }
func (s PutRequest) Timestamp() int64            { return int64(C.Struct(s).Get64(8)) }
func (s PutRequest) SetTimestamp(v int64)        { C.Struct(s).Set64(8, uint64(v)) }
func (s PutRequest) IndexVals() C.TextList       { return C.TextList(C.Struct(s).GetObject(0)) }
func (s PutRequest) SetIndexVals(v C.TextList)   { C.Struct(s).SetObject(0, C.Object(v)) }
func (s PutRequest) Payload() []byte             { return C.Struct(s).GetObject(1).ToData() }
func (s PutRequest) SetPayload(v []byte)         { C.Struct(s).SetObject(1, s.Segment.NewData(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s PutRequest) MarshalJSON() (bs []byte, err error) { return }

type PutRequest_List C.PointerList

func NewPutRequestList(s *C.Segment, sz int) PutRequest_List {
	return PutRequest_List(s.NewCompositeList(16, 2, sz))
}
func (s PutRequest_List) Len() int            { return C.PointerList(s).Len() }
func (s PutRequest_List) At(i int) PutRequest { return PutRequest(C.PointerList(s).At(i).ToStruct()) }
func (s PutRequest_List) ToArray() []PutRequest {
	n := s.Len()
	a := make([]PutRequest, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s PutRequest_List) Set(i int, item PutRequest) { C.PointerList(s).Set(i, C.Object(item)) }
