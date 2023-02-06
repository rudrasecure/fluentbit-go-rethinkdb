package main

import (
	"C"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"
)
import (
	"encoding/json"
	"log"

	"github.com/rudrasecure/fluentbit-go-rethinkdb/db"
)

var pluginName = "fluentbit-go-rethinkdb"
var r *db.RethinkDB

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
	database := output.FLBPluginConfigKey(plugin, "Database")
	tableName := output.FLBPluginConfigKey(plugin, "TableName")
	logKey := output.FLBPluginConfigKey(plugin, "LogKey")

	r = &db.RethinkDB{}

	err := r.Connect(connectionUri, database, tableName)
	if err != nil {
		log.Printf("[%s] Error connecting to RethinkDB: %s", pluginName, err)
		return output.FLB_ERROR
	}

	output.FLBPluginSetContext(plugin, logKey)

	return output.FLB_OK
}

//export FLBPluginFlushCtx
func FLBPluginFlushCtx(ctx, data unsafe.Pointer, length C.int, tag *C.char) int {
	decoder := output.NewDecoder(data, int(length))
	var logRecords []map[string]any

	for {
		ret, _, record := output.GetRecord(decoder)
		if ret != 0 {
			break
		}

		logLine := make(map[string]any)

		logKey := output.FLBPluginGetContext(ctx).(string)

		log.Printf("[%s] LogKey: %s", pluginName, logKey)

		switch record[logKey].(type) {
			case string:
				err := json.Unmarshal([]byte(record[logKey].(string)), &logLine)
				if err != nil {
					log.Printf("[%s] Error unmarshalling log: %s", pluginName, err)
				}

			case []uint8:
				err := json.Unmarshal(record[logKey].([]uint8), &logLine)
				if err != nil {
					log.Printf("[%s] Error unmarshalling log: %s", pluginName, err)
				}
		}

		logRecords = append(logRecords, logLine)
	}

	err := r.Insert(logRecords)
	if err != nil {
		log.Printf("[%s] Error inserting data to RethinkDB: %s", pluginName, err)
		return output.FLB_RETRY
	}

	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	err := r.Close()
	if err != nil {
		log.Printf("[%s] Error closing connection to RethinkDB: %s", pluginName, err)
	}
	return output.FLB_OK
}

func main() {
}