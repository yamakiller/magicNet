package library

import (
	"database/sql"
	"log"
	"time"

	// import mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// MySQLDB : 关系数据库
type MySQLDB struct {
	dsn string
	db  *sql.DB
}

// Init : 初始化mysql数据库
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

// Query : 查询
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
				record[columns[i]] = col
			}
		}
	}
	return record, nil
}

// Insert : 插入语句
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

// Update : 更新语句
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

// Close : 关闭
func (slf *MySQLDB) Close() {
	slf.db.Close()
}
