package library

import (
	"database/sql"
	"log"
	"reflect"
	"time"

	// import mysql driver
	_ "github.com/go-sql-driver/mysql"
)

type MySQLValue struct {
	vtype reflect.Type
	v     interface{}
}

func (slf *MySQLValue) ToString() string {
	return string(slf.v.([]uint8))
}

func (slf *MySQLValue) ToInt64() int64 {
	return slf.v.(int64)
}

//MySQLDB
//@struct MySQLDB desc: mysql operation object
//@member (string) mysql connection dsn
//@member (*sql.DB) mysql connection object
type MySQLDB struct {
	dsn string
	db  *sql.DB
}

//Init
//@method Init desc: initialization mysql DB
//@param (string) mysql connection dsn
//@param (int) mysql connection max of number
//@param (int) mysql connection idle of number
//@param (int) mysql connection life time[util/sec]
//@return (error) fail:return error, success: return nil
func (slf *MySQLDB) Init(dsn string, maxConn int, maxIdleConn, lifeSec int) error {
	slf.db, _ = sql.Open("mysql", dsn)
	slf.db.SetMaxOpenConns(maxConn)
	slf.db.SetMaxIdleConns(maxIdleConn)
	slf.db.SetConnMaxLifetime(time.Duration(lifeSec) * time.Second)
	slf.dsn = dsn

	err := slf.db.Ping()
	if err != nil {
		return err
	}

	return nil
}

//Query
//@method Query desc: execute sql query
//@param (string) query sql
//@param (...interface{}) sql params
//@return (map[string]interface{}) query result
//@return (error) fail: return error, success: return nil
func (slf *MySQLDB) Query(strSQL string, args ...interface{}) (map[string]interface{}, error) {
	if perr := slf.db.Ping(); perr != nil {
		return nil, perr
	}

	rows, err := slf.db.Query(strSQL, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]interface{})
	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = &MySQLValue{vtype: reflect.TypeOf(col), v: col}
			}
		}
	}
	return record, nil
}

//Insert
//@method Insert desc: execute sql Insert
//@param (string) Insert sql
//@param (...interface{}) sql params
//@return (int54) insert of number
//@return (error) fail: return error, success: return nil
func (slf *MySQLDB) Insert(strSQL string, args ...interface{}) (int64, error) {
	if perr := slf.db.Ping(); perr != nil {
		log.Println(perr)
		return 0, perr
	}

	r, err := slf.db.Exec(strSQL, args...)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return r.LastInsertId()
}

//Update
//@method Update desc: execute sql Update
//@param (string) Update sql
//@param (...interface{}) sql params
//@return (int54) Update of number
func (slf *MySQLDB) Update(strSQL string, args ...interface{}) (int64, error) {
	if perr := slf.db.Ping(); perr != nil {
		return 0, perr
	}

	r, err := slf.db.Exec(strSQL, args)
	if err != nil {
		return 0, err
	}

	return r.RowsAffected()
}

//Close
//@method CLose desc: close mysql connection
func (slf *MySQLDB) Close() {
	slf.db.Close()
}
