using Go = import "../../../glycerine/go-capnproto/go.capnp";

@0x97296a4ab461f129;
$Go.package("kmdb");
$Go.import("github.com/glycerine/go-capnproto/capnpc-go");

using Id = UInt64;

struct PutRequest {
  db @0 :Text;
  time @1 :Int64;
  values @2 :List(Text);
  payload @3 :Data;
}

struct PutResult {
  ok @0 :Bool;
}

struct GetRequest {
  db @0 :Text;
  start @1 :Int64;
  end @2 :Int64;
  values @3 :List(Text);
}

struct GetResult {
  ok @0 :Bool;
  data @1 :List(ResultItem);
}

struct ResultItem {
  values @0 :List(Text);
  data @1 :List(Data);
}
