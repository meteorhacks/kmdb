package kmdb

// AUTO GENERATED - DO NOT EDIT

import (
	C "github.com/glycerine/go-capnproto"
)

type PutRequest C.Struct

func NewPutRequest(s *C.Segment) PutRequest      { return PutRequest(s.NewStruct(8, 2)) }
func NewRootPutRequest(s *C.Segment) PutRequest  { return PutRequest(s.NewRootStruct(8, 2)) }
func AutoNewPutRequest(s *C.Segment) PutRequest  { return PutRequest(s.NewStructAR(8, 2)) }
func ReadRootPutRequest(s *C.Segment) PutRequest { return PutRequest(s.Root(0).ToStruct()) }
func (s PutRequest) Timestamp() int64            { return int64(C.Struct(s).Get64(0)) }
func (s PutRequest) SetTimestamp(v int64)        { C.Struct(s).Set64(0, uint64(v)) }
func (s PutRequest) IndexVals() C.TextList       { return C.TextList(C.Struct(s).GetObject(0)) }
func (s PutRequest) SetIndexVals(v C.TextList)   { C.Struct(s).SetObject(0, C.Object(v)) }
func (s PutRequest) Payload() []byte             { return C.Struct(s).GetObject(1).ToData() }
func (s PutRequest) SetPayload(v []byte)         { C.Struct(s).SetObject(1, s.Segment.NewData(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s PutRequest) MarshalJSON() (bs []byte, err error) { return }

type PutRequest_List C.PointerList

func NewPutRequestList(s *C.Segment, sz int) PutRequest_List {
	return PutRequest_List(s.NewCompositeList(8, 2, sz))
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

type PutResult C.Struct

func NewPutResult(s *C.Segment) PutResult      { return PutResult(s.NewStruct(8, 0)) }
func NewRootPutResult(s *C.Segment) PutResult  { return PutResult(s.NewRootStruct(8, 0)) }
func AutoNewPutResult(s *C.Segment) PutResult  { return PutResult(s.NewStructAR(8, 0)) }
func ReadRootPutResult(s *C.Segment) PutResult { return PutResult(s.Root(0).ToStruct()) }
func (s PutResult) Ok() bool                   { return C.Struct(s).Get1(0) }
func (s PutResult) SetOk(v bool)               { C.Struct(s).Set1(0, v) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s PutResult) MarshalJSON() (bs []byte, err error) { return }

type PutResult_List C.PointerList

func NewPutResultList(s *C.Segment, sz int) PutResult_List {
	return PutResult_List(s.NewCompositeList(8, 0, sz))
}
func (s PutResult_List) Len() int           { return C.PointerList(s).Len() }
func (s PutResult_List) At(i int) PutResult { return PutResult(C.PointerList(s).At(i).ToStruct()) }
func (s PutResult_List) ToArray() []PutResult {
	n := s.Len()
	a := make([]PutResult, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s PutResult_List) Set(i int, item PutResult) { C.PointerList(s).Set(i, C.Object(item)) }
