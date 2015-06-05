package proto

// AUTO GENERATED - DO NOT EDIT

import (
	C "github.com/glycerine/go-capnproto"
)

type PutRequest C.Struct

func NewPutRequest(s *C.Segment) PutRequest      { return PutRequest(s.NewStruct(8, 2)) }
func NewRootPutRequest(s *C.Segment) PutRequest  { return PutRequest(s.NewRootStruct(8, 2)) }
func AutoNewPutRequest(s *C.Segment) PutRequest  { return PutRequest(s.NewStructAR(8, 2)) }
func ReadRootPutRequest(s *C.Segment) PutRequest { return PutRequest(s.Root(0).ToStruct()) }
func (s PutRequest) Time() int64                 { return int64(C.Struct(s).Get64(0)) }
func (s PutRequest) SetTime(v int64)             { C.Struct(s).Set64(0, uint64(v)) }
func (s PutRequest) Values() C.TextList          { return C.TextList(C.Struct(s).GetObject(0)) }
func (s PutRequest) SetValues(v C.TextList)      { C.Struct(s).SetObject(0, C.Object(v)) }
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

type GetRequest C.Struct

func NewGetRequest(s *C.Segment) GetRequest      { return GetRequest(s.NewStruct(16, 1)) }
func NewRootGetRequest(s *C.Segment) GetRequest  { return GetRequest(s.NewRootStruct(16, 1)) }
func AutoNewGetRequest(s *C.Segment) GetRequest  { return GetRequest(s.NewStructAR(16, 1)) }
func ReadRootGetRequest(s *C.Segment) GetRequest { return GetRequest(s.Root(0).ToStruct()) }
func (s GetRequest) Start() int64                { return int64(C.Struct(s).Get64(0)) }
func (s GetRequest) SetStart(v int64)            { C.Struct(s).Set64(0, uint64(v)) }
func (s GetRequest) End() int64                  { return int64(C.Struct(s).Get64(8)) }
func (s GetRequest) SetEnd(v int64)              { C.Struct(s).Set64(8, uint64(v)) }
func (s GetRequest) Values() C.TextList          { return C.TextList(C.Struct(s).GetObject(0)) }
func (s GetRequest) SetValues(v C.TextList)      { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s GetRequest) MarshalJSON() (bs []byte, err error) { return }

type GetRequest_List C.PointerList

func NewGetRequestList(s *C.Segment, sz int) GetRequest_List {
	return GetRequest_List(s.NewCompositeList(16, 1, sz))
}
func (s GetRequest_List) Len() int            { return C.PointerList(s).Len() }
func (s GetRequest_List) At(i int) GetRequest { return GetRequest(C.PointerList(s).At(i).ToStruct()) }
func (s GetRequest_List) ToArray() []GetRequest {
	n := s.Len()
	a := make([]GetRequest, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s GetRequest_List) Set(i int, item GetRequest) { C.PointerList(s).Set(i, C.Object(item)) }

type GetResult C.Struct

func NewGetResult(s *C.Segment) GetResult      { return GetResult(s.NewStruct(8, 1)) }
func NewRootGetResult(s *C.Segment) GetResult  { return GetResult(s.NewRootStruct(8, 1)) }
func AutoNewGetResult(s *C.Segment) GetResult  { return GetResult(s.NewStructAR(8, 1)) }
func ReadRootGetResult(s *C.Segment) GetResult { return GetResult(s.Root(0).ToStruct()) }
func (s GetResult) Ok() bool                   { return C.Struct(s).Get1(0) }
func (s GetResult) SetOk(v bool)               { C.Struct(s).Set1(0, v) }
func (s GetResult) Data() ResultItem_List      { return ResultItem_List(C.Struct(s).GetObject(0)) }
func (s GetResult) SetData(v ResultItem_List)  { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s GetResult) MarshalJSON() (bs []byte, err error) { return }

type GetResult_List C.PointerList

func NewGetResultList(s *C.Segment, sz int) GetResult_List {
	return GetResult_List(s.NewCompositeList(8, 1, sz))
}
func (s GetResult_List) Len() int           { return C.PointerList(s).Len() }
func (s GetResult_List) At(i int) GetResult { return GetResult(C.PointerList(s).At(i).ToStruct()) }
func (s GetResult_List) ToArray() []GetResult {
	n := s.Len()
	a := make([]GetResult, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s GetResult_List) Set(i int, item GetResult) { C.PointerList(s).Set(i, C.Object(item)) }

type ResultItem C.Struct

func NewResultItem(s *C.Segment) ResultItem      { return ResultItem(s.NewStruct(0, 2)) }
func NewRootResultItem(s *C.Segment) ResultItem  { return ResultItem(s.NewRootStruct(0, 2)) }
func AutoNewResultItem(s *C.Segment) ResultItem  { return ResultItem(s.NewStructAR(0, 2)) }
func ReadRootResultItem(s *C.Segment) ResultItem { return ResultItem(s.Root(0).ToStruct()) }
func (s ResultItem) Values() C.TextList          { return C.TextList(C.Struct(s).GetObject(0)) }
func (s ResultItem) SetValues(v C.TextList)      { C.Struct(s).SetObject(0, C.Object(v)) }
func (s ResultItem) Data() C.DataList            { return C.DataList(C.Struct(s).GetObject(1)) }
func (s ResultItem) SetData(v C.DataList)        { C.Struct(s).SetObject(1, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s ResultItem) MarshalJSON() (bs []byte, err error) { return }

type ResultItem_List C.PointerList

func NewResultItemList(s *C.Segment, sz int) ResultItem_List {
	return ResultItem_List(s.NewCompositeList(0, 2, sz))
}
func (s ResultItem_List) Len() int            { return C.PointerList(s).Len() }
func (s ResultItem_List) At(i int) ResultItem { return ResultItem(C.PointerList(s).At(i).ToStruct()) }
func (s ResultItem_List) ToArray() []ResultItem {
	n := s.Len()
	a := make([]ResultItem, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s ResultItem_List) Set(i int, item ResultItem) { C.PointerList(s).Set(i, C.Object(item)) }
