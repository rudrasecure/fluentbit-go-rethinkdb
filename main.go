package main

import (
	"C"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"
)
import (
	"encoding/json"
	"log"
	"time"

	"github.com/rudrasecure/fluentbit-go-rethinkdb/db"
)

var pluginName = "fluentbit-go-rethinkdb"

//export FLBPluginRegister
func FLBPluginRegister(plugin unsafe.Pointer) int {
	return output.FLBPluginRegister(plugin, pluginName, "A Fluent Bit Go plugin for RethinkDB.")
}

//export FLBPluginInit
// (fluentbit will call this)
// plugin (context) pointer to fluentbit context (state/ c code)
func FLBPluginInit(plugin unsafe.Pointer) int {
	
	log.Printf("[%s] Init called", pluginName)
	connectionUri := output.FLBPluginConfigKey(plugin, "ConnectionUri")

	if connectionUri == "" {
		log.Printf("[%s] ConnectionUri is required", pluginName)
		return output.FLB_ERROR
	}

	database := output.FLBPluginConfigKey(plugin, "Database")
	tableName := output.FLBPluginConfigKey(plugin, "TableName")
	primaryKey := output.FLBPluginConfigKey(plugin, "PrimaryKey")

	if primaryKey == "" {
		primaryKey = "id"
	}

	if database == "" {
		database = "test"
	}

	if tableName == "" {
		tableName = "logs"
	}

	r := &db.RethinkDB{}

	err := r.Connect(connectionUri, database, tableName, primaryKey)
	if err != nil {
		log.Printf("[%s] Error connecting to RethinkDB: %s", pluginName, err)
		return output.FLB_ERROR
	}

	output.FLBPluginSetContext(plugin, r)

	return output.FLB_OK
}

//export FLBPluginFlushCtx
func FLBPluginFlushCtx(ctx, data unsafe.Pointer, length C.int, tag *C.char) int {
	decoder := output.NewDecoder(data, int(length))
	var logRecords []map[string]any
	r := output.FLBPluginGetContext(ctx).(*db.RethinkDB)

	for {
		ret, ts, record := output.GetRecord(decoder)
		if ret != 0 {
			break
		}

		var timeStamp time.Time
		switch t := ts.(type) {
		case output.FLBTime:
			timeStamp = ts.(output.FLBTime).Time
		case uint64:
			timeStamp = time.Unix(int64(t), 0)
		default:
			log.Print("given time is not in a known format, defaulting to now.\n")
			timeStamp = time.Now()
		}

		logLine := createJSON(&timeStamp, record)

		logRecords = append(logRecords, logLine)
	}

	err := r.Insert(logRecords)
	if err != nil {
		log.Printf("[%s] Error inserting data to RethinkDB: %s", pluginName, err)
		return output.FLB_RETRY
	}

	return output.FLB_OK
}

//export FLBPluginExitCtx
func FLBPluginExitCtx(ctx unsafe.Pointer) int {
	log.Printf("[%v] Exit called", ctx)
	r, ok := output.FLBPluginGetContext(ctx).(*db.RethinkDB)
	if ok {
		err := r.Close()
		if err != nil {
			log.Printf("[%s] Error closing connection to RethinkDB: %s", pluginName, err)
			return output.FLB_ERROR
		}
	}
	return output.FLB_OK
}

func parseMap(timestamp *time.Time, mapInterface map[interface{}]interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	
	for k, v := range mapInterface {
		switch t := v.(type) {
		case []byte:
			var data map[string]interface{}
			err := json.Unmarshal(t, &data)
			if err != nil {
				m[k.(string)] = string(t)
				continue
			}
			m[k.(string)] = data
		case map[interface{}]interface{}:
			m[k.(string)] = parseMap(nil, t)
		default:
			m[k.(string)] = v
		}
	}

	m["fluentbit_timestamp"] = timestamp
	return m
}

func createJSON(timestamp *time.Time, record map[interface{}]interface{}) map[string]interface{} {
	m := parseMap(timestamp, record)
	return m
}

func handlePanic() {
	if r := recover(); r != nil {
		log.Printf("[%s] Recovered in f %v", pluginName, r)
	}
}

func main() {
	defer handlePanic()
}