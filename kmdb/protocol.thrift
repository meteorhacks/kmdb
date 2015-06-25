namespace go kmdb


//   Errors
// ----------

typedef i64 ErrorCode
const ErrorCode ERR_NO_ERROR = 0
const ErrorCode ERR_DB_NOT_FOUND = 1
const ErrorCode ERR_BATCH_ERROR = 1


//   Types
// ---------

struct PutReq {
  1: string database,
  2: i64 timestamp,
  3: list<string> fields,
  4: double value,
  5: i64 count,
}

struct PutRes {
  1: ErrorCode error = ERR_NO_ERROR,
}

struct IncReq {
  1: string database,
  2: i64 timestamp,
  3: list<string> fields,
  4: double value,
  5: i64 count,
}

struct IncRes {
  1: ErrorCode error = ERR_NO_ERROR,
}

struct GetReq {
  1: string database,
  2: i64 startTime,
  3: i64 endTime,
  4: list<string> fields,
  5: list<bool> groupBy,
}

struct GetRes {
  1: ErrorCode error = ERR_NO_ERROR,
  2: list<ResSeries> data,
}

struct ResSeries {
  1: list<string> fields,
  2: list<ResPoint> points,
}

struct ResPoint {
  1: double value,
  2: i64 count,
}


//   Services
// ------------

service ThriftService {
  PutRes put(1: PutReq req),
  IncRes inc(1: IncReq req),
  GetRes get(1: GetReq req),
  list<PutRes> putBatch(1: list<PutReq> batch),
  list<IncRes> incBatch(1: list<IncReq> batch),
  list<GetRes> getBatch(1: list<GetReq> batch),
}
