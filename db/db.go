package db

import (
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type RethinkDB struct{
	session *r.Session
	tableName string
}

func (rdb *RethinkDB) Connect(connectionUri string, database string, tableName string) error {
	session, err := r.Connect(r.ConnectOpts {
		Address:  connectionUri,
		Database: database,
	})

	rdb.session = session
	rdb.tableName = tableName
	if err != nil {
		return err
	}
	return nil
}

func (rdb *RethinkDB) Insert(data interface{}) error {
	_, err := r.Table(rdb.tableName).Insert(data).RunWrite(rdb.session)
	if err != nil {
		return err
	}
	return nil
}
