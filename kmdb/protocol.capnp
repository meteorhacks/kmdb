using Go = import "../../../glycerine/go-capnproto/go.capnp";

@0x97296a4ab461f129;
$Go.package("kmdb");
$Go.import("github.com/glycerine/go-capnproto/capnpc-go");

using Id = UInt64;

#   Put
# -------

struct PutRequest {
  database @0 :Text;
  timestamp @1 :Int64;
  fields @2 :List(Text);
  value @3 :Float64;
  count @4 :Int64;
}

struct PutResult {
  ok @0 :Bool;
}

#   Get
# -------

struct GetRequest {
  database @0 :Text;
  startTime @1 :Int64;
  endTime @2 :Int64;
  fields @3 :List(Text);
  groupBy @4 :List(Bool);
}

struct GetResult {
  ok @0 :Bool;
  data @1 :List(ResultSeries);
}

struct ResultSeries {
  fields @0 :List(Text);
  points @1 :List(ResultPoint);
}

struct ResultPoint {
  value @0 :Float64;
  count @1 :Int64;
}
