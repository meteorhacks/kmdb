package kmdb

type Config struct {
	// database name. Currently only used with naming files
	// can be useful when supporting multiple Databases
	DatabaseName string `json:"databaseName"`

	// place to store data files
	DataPath string `json:"dataPath"`

	// depth of the index tree
	IndexDepth int64 `json:"indexDepth"`

	// payload size should always be equal to this amount
	PayloadSize int64 `json:"payloadSize"`

	// time duration in nano seconds of a range unit
	// this should be a multiple of `Resolution`
	BucketDuration int64 `json:"bucketDuration"`

	// bucket resolution in nano seconds
	Resolution int64 `json:"resolution"`

	// number of records per segment
	SegmentSize int64 `json:"segmentSize"`

	// enable pprof on ":6060" instead of "localhost:6060".
	DebugMode bool `json:"debugMode"`

	// address to listen for ddp traffic (host:port)
	BDDPAddress string `json:"bddpAddress"`
}
