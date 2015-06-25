# KMDB
A database only good for storing metric data. KMDB is powered by [kdb](https://github.com/meteorhacks/kdb), a lightweight database engine which is optimized for high write throughput.


## Configuration

Create a configuration JSON file with following fields. All fields are mandatory. When starting the server, use `kmdb -config /path/to/config.json` to use your settings.

```json
{
    "remote_debug": true,
    "listen_address": ":3000",

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

## Docker

KMDB uses memory mapping to increase write performance therefore for KMDB to work the `IPC_LOCK` linux capability must be enabled when running inside docker. This can be done easily by adding `--cap-add=IPC_LOCK` when starting the container.

```bash
docker run -d \
  -p 3000:3000 \
  -v /data:/data \
  -v /etc/kmdb.json:/etc/kmdb.json \
  --cap-add=IPC_LOCK \
  meteorhacks/kmdb
```

***

**Note: KMDB is under active development therefor not suitable for production use.**
