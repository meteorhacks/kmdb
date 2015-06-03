using Go = import "../../glycerine/go-capnproto/go.capnp";

@0x97296a4ab461f129;
$Go.package("kmdb");
$Go.import("github.com/glycerine/go-capnproto/capnpc-go");

using Id = UInt64;

struct PutRequest {
  timestamp @0 :Int64;
  indexVals @1 :List(Text);
  payload @2 :Data;
}

struct PutResult {
  ok @0 :Bool;
}
