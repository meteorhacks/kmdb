using Go = import "../../../glycerine/go-capnproto/go.capnp";

@0x97296a4ab461f129;
$Go.package("proto");
$Go.import("github.com/glycerine/go-capnproto/capnpc-go");

using Id = UInt64;

struct PutRequest {
  time @0 :Int64;
  values @1 :List(Text);
  payload @2 :Data;
}

struct PutResult {
  ok @0 :Bool;
}

struct GetRequest {
  start @0 :Int64;
  end @1 :Int64;
  values @2 :List(Text);
}

struct GetResult {
  ok @0 :Bool;
  data @1 :List(ResultItem);
}

struct ResultItem {
  values @0 :List(Text);
  data @1 :List(Data);
}
