# fluentbit-go-rethinkdb
Fluent Bit RethinkDB output plugin written in Golang

## Build
* You need to first compile Fluent Bit from source with Golang support as mentioned [here](https://docs.fluentbit.io/manual/development/golang-output-plugins).
It is quite possible that the binary you have installed already has support for accepting a binary plugin as a parameter. You can check that by running the following:
```
$ fluent-bit -h
Usage: fluent-bit [OPTION]

Available Options
  -b  --storage_path=PATH specify a storage buffering path
  -c  --config=FILE       specify an optional configuration file
  -d, --daemon            run Fluent Bit in background mode
  -D, --dry-run           dry run
  -f, --flush=SECONDS     flush timeout in seconds (default: 1)
  -C, --custom=CUSTOM     enable a custom plugin
  -i, --input=INPUT       set an input
  -F  --filter=FILTER     set a filter
  -m, --match=MATCH       set plugin match, same as '-p match=abc'
  -o, --output=OUTPUT     set an output
  -p, --prop="A=B"        set plugin configuration property
  -R, --parser=FILE       specify a parser configuration file
  -e, --plugin=FILE       load an external plugin (shared lib)
  -l, --log_file=FILE     write log info to a file
  -t, --tag=TAG           set plugin tag, same as '-p tag=abc'
  -T, --sp-task=SQL       define a stream processor task
  -v, --verbose           increase logging verbosity (default: info)
  -Z, --enable-chunk-traceenable chunk tracing. activating it requires using the HTTP Server API.
  -w, --workdir           set the working directory
  -H, --http              enable monitoring HTTP server
  -P, --port              set HTTP server TCP port (default: 2020)
  -s, --coro_stack_size   set coroutines stack size in bytes (default: 12288)
  -q, --quiet             quiet mode
  -S, --sosreport         support report for Enterprise customers
  -V, --version           show version number
  -h, --help              print this help
```
You should have an `-e` option available to load an external plugin as a binary.
* In the root of the directory, run
```
go build -buildmode=c-shared -o out/fluentbit-go-rethinkdb.so
```
This will generate the binary for the plugin that can now be loaded while running `fluent-bit` using
```
fluent-bit --plugin=./out/fluentbit-go-rethinkdb.so -c fluent-bit.conf
```

## Configuration
The project contains an example `fluent-bit.conf` file explaining the configuration supported by the plugin
```
[OUTPUT]
    Name          fluentbit-go-rethinkdb  # Name of the plugin. DO NOT CHANGE!
    Match         *                       # This will read everything coming from the input without any filters
    ConnectionUri localhost:28015         # RethinkDB connection URL
    Database      logs                    # RethinkDB DB name
    TableName     cpu                     # RethinkDB Table name
    LogKey        log                     # Each record is read as a map[string]any and this key is used to read the log data from the map
```

## Example

* Input log line after running tail
```json
{"pkts_toserver":1,"pkts_toclient":0,"bytes_toserver":54,"bytes_toclient":0,"start":"2023-02-04T16:27:06.058562+0000","end":"2023-02-04T16:27:06.058562+0000","age":0,"state":"new","reason":"timeout","alerted":false,"community_id":"1:D2i50Uk5MR+GQETnkjI9zP8nvOc=","tcp":{"tcp_flags":"00","tcp_flaDDDDDDDDDDDDDDDDgs_ts":"00","tcp_flags_tc":"004"},"host":"rumars-dep-asasrudra-3dadf162sdsd-0e31-4997-bb58-c7bf74c570dsddswe8"}

```

* Output in RethinkDB Admin Dashboard
![image](https://user-images.githubusercontent.com/5111523/216781658-cd8acd34-9d22-425e-b3ad-e170439fbefa.png)
