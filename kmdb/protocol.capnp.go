package kmdb

// AUTO GENERATED - DO NOT EDIT

import (
	C "github.com/glycerine/go-capnproto"
	"math"
)

type PutRequest C.Struct

func NewPutRequest(s *C.Segment) PutRequest      { return PutRequest(s.NewStruct(24, 2)) }
func NewRootPutRequest(s *C.Segment) PutRequest  { return PutRequest(s.NewRootStruct(24, 2)) }
func AutoNewPutRequest(s *C.Segment) PutRequest  { return PutRequest(s.NewStructAR(24, 2)) }
func ReadRootPutRequest(s *C.Segment) PutRequest { return PutRequest(s.Root(0).ToStruct()) }
func (s PutRequest) Database() string            { return C.Struct(s).GetObject(0).ToText() }
func (s PutRequest) SetDatabase(v string)        { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s PutRequest) Timestamp() int64            { return int64(C.Struct(s).Get64(0)) }
func (s PutRequest) SetTimestamp(v int64)        { C.Struct(s).Set64(0, uint64(v)) }
func (s PutRequest) Fields() C.TextList          { return C.TextList(C.Struct(s).GetObject(1)) }
func (s PutRequest) SetFields(v C.TextList)      { C.Struct(s).SetObject(1, C.Object(v)) }
func (s PutRequest) Value() float64              { return math.Float64frombits(C.Struct(s).Get64(8)) }
func (s PutRequest) SetValue(v float64)          { C.Struct(s).Set64(8, math.Float64bits(v)) }
func (s PutRequest) Count() int64                { return int64(C.Struct(s).Get64(16)) }
func (s PutRequest) SetCount(v int64)            { C.Struct(s).Set64(16, uint64(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s PutRequest) MarshalJSON() (bs []byte, err error) { return }

type PutRequest_List C.PointerList

func NewPutRequestList(s *C.Segment, sz int) PutRequest_List {
	return PutRequest_List(s.NewCompositeList(24, 2, sz))
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

type IncRequest C.Struct

func NewIncRequest(s *C.Segment) IncRequest      { return IncRequest(s.NewStruct(24, 2)) }
func NewRootIncRequest(s *C.Segment) IncRequest  { return IncRequest(s.NewRootStruct(24, 2)) }
func AutoNewIncRequest(s *C.Segment) IncRequest  { return IncRequest(s.NewStructAR(24, 2)) }
func ReadRootIncRequest(s *C.Segment) IncRequest { return IncRequest(s.Root(0).ToStruct()) }
func (s IncRequest) Database() string            { return C.Struct(s).GetObject(0).ToText() }
func (s IncRequest) SetDatabase(v string)        { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s IncRequest) Timestamp() int64            { return int64(C.Struct(s).Get64(0)) }
func (s IncRequest) SetTimestamp(v int64)        { C.Struct(s).Set64(0, uint64(v)) }
func (s IncRequest) Fields() C.TextList          { return C.TextList(C.Struct(s).GetObject(1)) }
func (s IncRequest) SetFields(v C.TextList)      { C.Struct(s).SetObject(1, C.Object(v)) }
func (s IncRequest) Value() float64              { return math.Float64frombits(C.Struct(s).Get64(8)) }
func (s IncRequest) SetValue(v float64)          { C.Struct(s).Set64(8, math.Float64bits(v)) }
func (s IncRequest) Count() int64                { return int64(C.Struct(s).Get64(16)) }
func (s IncRequest) SetCount(v int64)            { C.Struct(s).Set64(16, uint64(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s IncRequest) MarshalJSON() (bs []byte, err error) { return }

type IncRequest_List C.PointerList

func NewIncRequestList(s *C.Segment, sz int) IncRequest_List {
	return IncRequest_List(s.NewCompositeList(24, 2, sz))
}
func (s IncRequest_List) Len() int            { return C.PointerList(s).Len() }
func (s IncRequest_List) At(i int) IncRequest { return IncRequest(C.PointerList(s).At(i).ToStruct()) }
func (s IncRequest_List) ToArray() []IncRequest {
	n := s.Len()
	a := make([]IncRequest, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s IncRequest_List) Set(i int, item IncRequest) { C.PointerList(s).Set(i, C.Object(item)) }

type IncResult C.Struct

func NewIncResult(s *C.Segment) IncResult      { return IncResult(s.NewStruct(8, 0)) }
func NewRootIncResult(s *C.Segment) IncResult  { return IncResult(s.NewRootStruct(8, 0)) }
func AutoNewIncResult(s *C.Segment) IncResult  { return IncResult(s.NewStructAR(8, 0)) }
func ReadRootIncResult(s *C.Segment) IncResult { return IncResult(s.Root(0).ToStruct()) }
func (s IncResult) Ok() bool                   { return C.Struct(s).Get1(0) }
func (s IncResult) SetOk(v bool)               { C.Struct(s).Set1(0, v) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s IncResult) MarshalJSON() (bs []byte, err error) { return }

type IncResult_List C.PointerList

func NewIncResultList(s *C.Segment, sz int) IncResult_List {
	return IncResult_List(s.NewCompositeList(8, 0, sz))
}
func (s IncResult_List) Len() int           { return C.PointerList(s).Len() }
func (s IncResult_List) At(i int) IncResult { return IncResult(C.PointerList(s).At(i).ToStruct()) }
func (s IncResult_List) ToArray() []IncResult {
	n := s.Len()
	a := make([]IncResult, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s IncResult_List) Set(i int, item IncResult) { C.PointerList(s).Set(i, C.Object(item)) }

type GetRequest C.Struct

func NewGetRequest(s *C.Segment) GetRequest      { return GetRequest(s.NewStruct(16, 3)) }
func NewRootGetRequest(s *C.Segment) GetRequest  { return GetRequest(s.NewRootStruct(16, 3)) }
func AutoNewGetRequest(s *C.Segment) GetRequest  { return GetRequest(s.NewStructAR(16, 3)) }
func ReadRootGetRequest(s *C.Segment) GetRequest { return GetRequest(s.Root(0).ToStruct()) }
func (s GetRequest) Database() string            { return C.Struct(s).GetObject(0).ToText() }
func (s GetRequest) SetDatabase(v string)        { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s GetRequest) StartTime() int64            { return int64(C.Struct(s).Get64(0)) }
func (s GetRequest) SetStartTime(v int64)        { C.Struct(s).Set64(0, uint64(v)) }
func (s GetRequest) EndTime() int64              { return int64(C.Struct(s).Get64(8)) }
func (s GetRequest) SetEndTime(v int64)          { C.Struct(s).Set64(8, uint64(v)) }
func (s GetRequest) Fields() C.TextList          { return C.TextList(C.Struct(s).GetObject(1)) }
func (s GetRequest) SetFields(v C.TextList)      { C.Struct(s).SetObject(1, C.Object(v)) }
func (s GetRequest) GroupBy() C.BitList          { return C.BitList(C.Struct(s).GetObject(2)) }
func (s GetRequest) SetGroupBy(v C.BitList)      { C.Struct(s).SetObject(2, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s GetRequest) MarshalJSON() (bs []byte, err error) { return }

type GetRequest_List C.PointerList

func NewGetRequestList(s *C.Segment, sz int) GetRequest_List {
	return GetRequest_List(s.NewCompositeList(16, 3, sz))
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

func NewGetResult(s *C.Segment) GetResult       { return GetResult(s.NewStruct(8, 1)) }
func NewRootGetResult(s *C.Segment) GetResult   { return GetResult(s.NewRootStruct(8, 1)) }
func AutoNewGetResult(s *C.Segment) GetResult   { return GetResult(s.NewStructAR(8, 1)) }
func ReadRootGetResult(s *C.Segment) GetResult  { return GetResult(s.Root(0).ToStruct()) }
func (s GetResult) Ok() bool                    { return C.Struct(s).Get1(0) }
func (s GetResult) SetOk(v bool)                { C.Struct(s).Set1(0, v) }
func (s GetResult) Data() ResultSeries_List     { return ResultSeries_List(C.Struct(s).GetObject(0)) }
func (s GetResult) SetData(v ResultSeries_List) { C.Struct(s).SetObject(0, C.Object(v)) }

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

type ResultSeries C.Struct

func NewResultSeries(s *C.Segment) ResultSeries      { return ResultSeries(s.NewStruct(0, 2)) }
func NewRootResultSeries(s *C.Segment) ResultSeries  { return ResultSeries(s.NewRootStruct(0, 2)) }
func AutoNewResultSeries(s *C.Segment) ResultSeries  { return ResultSeries(s.NewStructAR(0, 2)) }
func ReadRootResultSeries(s *C.Segment) ResultSeries { return ResultSeries(s.Root(0).ToStruct()) }
func (s ResultSeries) Fields() C.TextList            { return C.TextList(C.Struct(s).GetObject(0)) }
func (s ResultSeries) SetFields(v C.TextList)        { C.Struct(s).SetObject(0, C.Object(v)) }
func (s ResultSeries) Points() ResultPoint_List      { return ResultPoint_List(C.Struct(s).GetObject(1)) }
func (s ResultSeries) SetPoints(v ResultPoint_List)  { C.Struct(s).SetObject(1, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s ResultSeries) MarshalJSON() (bs []byte, err error) { return }

type ResultSeries_List C.PointerList

func NewResultSeriesList(s *C.Segment, sz int) ResultSeries_List {
	return ResultSeries_List(s.NewCompositeList(0, 2, sz))
}
func (s ResultSeries_List) Len() int { return C.PointerList(s).Len() }
func (s ResultSeries_List) At(i int) ResultSeries {
	return ResultSeries(C.PointerList(s).At(i).ToStruct())
}
func (s ResultSeries_List) ToArray() []ResultSeries {
	n := s.Len()
	a := make([]ResultSeries, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s ResultSeries_List) Set(i int, item ResultSeries) { C.PointerList(s).Set(i, C.Object(item)) }

type ResultPoint C.Struct

func NewResultPoint(s *C.Segment) ResultPoint      { return ResultPoint(s.NewStruct(16, 0)) }
func NewRootResultPoint(s *C.Segment) ResultPoint  { return ResultPoint(s.NewRootStruct(16, 0)) }
func AutoNewResultPoint(s *C.Segment) ResultPoint  { return ResultPoint(s.NewStructAR(16, 0)) }
func ReadRootResultPoint(s *C.Segment) ResultPoint { return ResultPoint(s.Root(0).ToStruct()) }
func (s ResultPoint) Value() float64               { return math.Float64frombits(C.Struct(s).Get64(0)) }
func (s ResultPoint) SetValue(v float64)           { C.Struct(s).Set64(0, math.Float64bits(v)) }
func (s ResultPoint) Count() int64                 { return int64(C.Struct(s).Get64(8)) }
func (s ResultPoint) SetCount(v int64)             { C.Struct(s).Set64(8, uint64(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s ResultPoint) MarshalJSON() (bs []byte, err error) { return }

type ResultPoint_List C.PointerList

func NewResultPointList(s *C.Segment, sz int) ResultPoint_List {
	return ResultPoint_List(s.NewCompositeList(16, 0, sz))
}
func (s ResultPoint_List) Len() int             { return C.PointerList(s).Len() }
func (s ResultPoint_List) At(i int) ResultPoint { return ResultPoint(C.PointerList(s).At(i).ToStruct()) }
func (s ResultPoint_List) ToArray() []ResultPoint {
	n := s.Len()
	a := make([]ResultPoint, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s ResultPoint_List) Set(i int, item ResultPoint) { C.PointerList(s).Set(i, C.Object(item)) }
