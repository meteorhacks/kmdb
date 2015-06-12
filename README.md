# KMDB
A database only good for storing metric data. KMDB is powered by [kdb](https://github.com/meteorhacks/kdb), a lightweight database engine which is optimized for high write throughput.


## Configuration

Create a configuration JSON file with following fields. All fields are mandatory. When starting the server, use `kmdb -config /path/to/config.json` to use your settings. *Note: config file format will change significantly when multiple db support is implemented.*

```json
{
    "remote_debug": true,
    "bddp_address": ":3000",

    "databases": {
        "test": {
            "database_path": "/tmp/test-db",
            "index_depth": 4,
            "payload_size": 16,
            "payload_resolution": 60000000000,
            "bucket_duration": 3600000000000,
            "segment_size": 100000
        }
    }
}
```

***

**Note: KMDB is under active development therefor not suitable for production use.**
