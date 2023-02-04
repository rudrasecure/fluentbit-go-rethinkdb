package main

import (
	"log"
	"unsafe"
	"C"
	"github.com/fluent/fluent-bit-go/output"
	"github.com/rudrasecure/fluentbit-go-rethinkdb/db"
)

var pluginName string = "fluentbit-go-rethinkdb"

func FLBPluginRegister(def unsafe.Pointer) int {
	log.Printf("[%s] Register called", pluginName)
	return output.FLBPluginRegister(def, pluginName, "An output plugin for Fluent Bit to send data to RethinkDB")
}

func FLBPluginInit(plugin unsafe.Pointer) int {
	log.Printf("[%s] Init called", pluginName)
	connectionUri := output.FLBPluginConfigKey(plugin, "ConnectionUri")
	database := output.FLBPluginConfigKey(plugin, "Database")
	tableName := output.FLBPluginConfigKey(plugin, "TableName")

	r := &db.RethinkDB{}

	err := r.Connect(connectionUri, database, tableName)
	if err != nil {
		log.Printf("[%s] Error connecting to RethinkDB: %s", pluginName, err)
		return output.FLB_ERROR
	}

	output.FLBPluginSetContext(plugin, r)
	return output.FLB_OK
}

func FLBPluginFlush(data unsafe.Pointer, length C.int, tag *C.char) int {
	log.Printf("[%s] Flush called", pluginName)
	r := output.FLBPluginGetContext(data).(*db.RethinkDB)

	decoder := output.NewDecoder(data, int(length))
	var logRecords []map[any]any

	for {
		ret, _, record := output.GetRecord(decoder)
		if ret != 0 {
			break
		}

		logRecords = append(logRecords, record)
	}

	err := r.Insert(logRecords)
	if err != nil {
		log.Printf("[%s] Error inserting data to RethinkDB: %s", pluginName, err)
		return output.FLB_ERROR
	}

	return output.FLB_OK
}

func FLBPluginExit() int {
	log.Printf("[%s] [info] exit", pluginName)
	return output.FLB_OK
}

func main() {

}