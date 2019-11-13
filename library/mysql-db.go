package library

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	// import mysql driver
	_ "github.com/go-sql-driver/mysql"
)

//MySQLValue desc
//@struct MySQLValue desc mysql result value
type MySQLValue struct {
	v interface{}
}

//Print desc
//@method Print desc:
func (slf *MySQLValue) Print() {
	fmt.Printf("%+v, %d, %+v\n", reflect.TypeOf(slf.v), len(slf.v.([]uint8)), slf.v)
}

//IsEmpty desc
//@method IsEmpty desc: Return Is Empty
//@return (bool) emtpy: true  no empty:false
func (slf *MySQLValue) IsEmpty() bool {
	if slf.v == nil {
		return true
	}
	return false
}

//ToString desc
//@method ToString desc: Return string value
//@return (string) a string
func (slf *MySQLValue) ToString() string {
	return slf.getString()
}

//ToUint desc
//@method ToUint desc: Return uint value
//@return (uint) a value
func (slf *MySQLValue) ToUint() uint {
	v, e := strconv.Atoi(slf.getString())
	if e != nil {
		return 0
	}

	return uint(v)
}

//ToInt desc
//@method ToInt desc: Return int value
//@return (int) a value
func (slf *MySQLValue) ToInt() int {
	v, e := strconv.Atoi(slf.getString())
	if e != nil {
		return 0
	}

	return v
}

//ToUint32 desc
//@method ToUint32 desc: Return uint32 value
//@return (uint32) a value
func (slf *MySQLValue) ToUint32() uint32 {
	v, e := strconv.Atoi(slf.getString())
	if e != nil {
		return 0
	}

	return uint32(v)
}

//ToInt32 desc
//@method ToInt32 desc: Return int32 value
//@return (int32) a value
func (slf *MySQLValue) ToInt32() int32 {
	v, e := strconv.Atoi(slf.getString())
	if e != nil {
		return 0
	}

	return int32(v)
}

//ToUint64 desc
//@method ToUint64 desc: Return uint64 value
//@return (uint64) a value
func (slf *MySQLValue) ToUint64() uint64 {
	v, e := strconv.ParseInt(slf.getString(), 10, 64)
	if e != nil {
		return 0
	}
	return uint64(v)
}

//ToInt64 desc
//@method ToInt64 desc: Return int64 value
//@return (int64) a value
func (slf *MySQLValue) ToInt64() int64 {
	v, e := strconv.ParseInt(slf.getString(), 10, 64)
	if e != nil {
		return 0
	}
	return v
}

//ToFloat desc
//@method ToFloat desc: Return float32 value
//@return (float32) a value
func (slf *MySQLValue) ToFloat() float32 {
	v, e := strconv.ParseFloat(slf.getString(), 32)
	if e != nil {
		return 0.0
	}
	return float32(v)
}

//ToDouble desc
//@method ToDouble desc: Return float64 value
//@return (float64) a value
func (slf *MySQLValue) ToDouble() float64 {
	v, e := strconv.ParseFloat(slf.getString(), 64)
	if e != nil {
		return 0.0
	}
	return v
}

//ToByte desc
//@method ToByte desc: Return []byte value
//@return ([]byte) a value
func (slf *MySQLValue) ToByte() []byte {
	return ([]byte)(slf.v.([]uint8))
}

//ToTimeStamp desc
//@method ToTimeStamp desc: Return  time int64 value
//@return (int64) a value
func (slf *MySQLValue) ToTimeStamp() int64 {
	v := slf.ToDateTime()
	if v == nil {
		return 0
	}

	return v.Unix()
}

//ToDate desc
//@method ToDate desc: Return  time date value
//@return (*time.Time) a value
func (slf *MySQLValue) ToDate() *time.Time {
	v, e := time.Parse("2006-01-02", slf.getString())
	if e != nil {
		return nil
	}

	return &v
}

//ToDateTime desc
//@method ToDateTime desc: Return  time date time value
//@return (*time.Time) a value
func (slf *MySQLValue) ToDateTime() *time.Time {
	v, e := time.Parse("2006-01-02 15:04:05", slf.getString())
	if e != nil {
		return nil
	}

	return &v
}

func (slf *MySQLValue) getString() string {
	return string(slf.v.([]uint8))
}

//MySQLReader desc
//@struct MySQLReader desc: mysql reader
//@member (int) count row of number
//@member (int) read current row in index
//@member ([]string) columns name
//@member ([]MySQLValue) a mysql value
type MySQLReader struct {
	rows       int
	currentRow int
	columns    []string
	data       []MySQLValue
}

//GetAsNameValue desc
//@method GetAsNameValue desc: Return column name to value
//@return (*MySQLValue) mysql value
//@return (error) error informat
func (slf *MySQLReader) GetAsNameValue(name string) (*MySQLValue, error) {
	idx := slf.getNamePos(name)
	if idx == -1 {
		return nil, fmt.Errorf("mysql column %s is does not exist", name)
	}

	return slf.GetValue(idx)
}

//GetValue desc
//@method GetValue desc: Return column index to value
//@return (*MySQLValue) mysql value
//@return (error) error informat
func (slf *MySQLReader) GetValue(idx int) (*MySQLValue, error) {
	rpos := (slf.currentRow * len(slf.columns)) + idx
	if rpos >= len(slf.data) {
		return nil, fmt.Errorf("mysql column %d overload", idx)
	}

	return &slf.data[rpos], nil
}

//Next desc
//@method Next desc: Return Next row is success
//@return (bool) Next Success: true Next Fail:false
func (slf *MySQLReader) Next() bool {
	if (slf.currentRow + 1) >= slf.rows {
		return false
	}
	slf.currentRow++
	return true
}

//GetColumn desc
//@method GetColumn desc: Return Column of number
//@return (int) a number
func (slf *MySQLReader) GetColumn() int {
	return len(slf.columns)
}

//GetRow desc
//@method GetRow desc: Return row of number
//@return (int) a number
func (slf *MySQLReader) GetRow() int {
	return int(slf.rows)
}

//Rest desc
//@method Rest desc: go to first row
func (slf *MySQLReader) Rest() {
	slf.currentRow = -1
}

//Close desc
//@method Close desc: clear data
func (slf *MySQLReader) Close() {
	slf.columns = nil
	slf.data = nil
}

func (slf *MySQLReader) getNamePos(name string) int {
	for i, v := range slf.columns {
		if v == name {
			return i
		}
	}
	return -1
}

//MySQLDB desc
//@struct MySQLDB desc: mysql operation object
//@member (string) mysql connection dsn
//@member (*sql.DB) mysql connection object
type MySQLDB struct {
	dsn string
	db  *sql.DB
}

//Init desc
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

//Query desc
//@method Query desc: execute sql query
//@param (string) query sql
//@param (...interface{}) sql params
//@return (map[string]interface{}) query result
//@return (error) fail: return error, success: return nil
func (slf *MySQLDB) Query(strSQL string, args ...interface{}) (*MySQLReader, error) {
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

	record := &MySQLReader{currentRow: -1}
	record.columns = columns
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		d := make([]MySQLValue, len(columns))
		for i, col := range values {
			d[i].v = col
		}
		record.data = append(record.data, d...)
		record.rows++
	}

	return record, nil
}

//Insert desc
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

//Update desc
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

//Close desc
//@method CLose desc: close mysql connection
func (slf *MySQLDB) Close() {
	slf.db.Close()
}
