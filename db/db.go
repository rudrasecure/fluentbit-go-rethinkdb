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
	if err != nil {
		return err
	}

	res, err := r.DBList().Contains(database).Run(session)
	if err != nil {
		return err
	}

	var exists bool
	err = res.One(&exists)

	if err != nil {
		return err
	}

	if !exists {
		_, err = r.DBCreate(database).RunWrite(session)
		if err != nil {
			return err
		}
	}

	res, err = r.TableList().Contains(tableName).Run(session)
	if err != nil {
		return err
	}

	err = res.One(&exists)
	if err != nil {
		return err
	}

	if !exists {
		_, err = r.DB(database).TableCreate(tableName).RunWrite(session)
		if err != nil {
			return err
		}
	}

	rdb.session = session
	rdb.tableName = tableName
	
	return nil
}

func (rdb *RethinkDB) Insert(data interface{}) error {
	_, err := r.Table(rdb.tableName).Insert(data).RunWrite(rdb.session)
	if err != nil {
		return err
	} 
	return nil
}

func (rdb *RethinkDB) Close() error {
	return rdb.session.Close()
}
