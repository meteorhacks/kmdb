# KMDB
A database only good for storing metric data. KMDB is powered by [kdb](https://github.com/meteorhacks/kdb), a lightweight database engine which is optimized for high write throughput.


## Configuration

Create a configuration JSON file with following fields. All fields are mandatory. When starting the server, use `kmdb -config /path/to/config.json` to use your settings. *Note: config file format will change significantly when multiple db support is implemented.*

```json
{
  "databaseName": "test",
  "dataPath": "/tmp/test-db",
  "indexDepth": 4,
  "payloadSize": 16,
  "bucketDuration": 3600000000000,
  "resolution": 60000000000,
  "segmentSize": 100000,

  "debugMode": true,
  "bddpAddress": ":3000"
}
```

***

**Note: KMDB is under active development therefor not suitable for production use.**
