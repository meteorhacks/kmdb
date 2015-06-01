using Go = import "../../glycerine/go-capnproto/go.capnp";

@0x97296a4ab461f129;
$Go.package("kmdb");
$Go.import("github.com/glycerine/go-capnproto/capnpc-go");

using Id = UInt64;

struct PutRequest {
  partition @0 :Int64;
  timestamp @1 :Int64;
  indexVals @2 :List(Text);
  payload @3 :Data;
}
