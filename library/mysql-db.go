package library

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLDB : 关系数据库
type MySQLDB struct {
	dsn string
	db  *sql.DB
}

// Init : 初始化mysql数据库
func (msd *MySQLDB) Init(dsn string, maxConn int, maxIdleConn, lifeSec int) error {
	msd.db, _ = sql.Open("mysql", dsn)
	msd.db.SetMaxOpenConns(maxConn)
	msd.db.SetMaxIdleConns(maxIdleConn)
	msd.db.SetConnMaxLifetime(time.Duration(lifeSec) * time.Second)
	msd.dsn = dsn

	err := msd.db.Ping()
	if err != nil {
		return err
	}

	return nil
}

// Query : 查询
func (msd *MySQLDB) Query(strSQL string, args ...interface{}) (map[string]interface{}, error) {
	if perr := msd.db.Ping(); perr != nil {
		return nil, perr
	}

	rows, err := msd.db.Query(strSQL, args)
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
func (msd *MySQLDB) Insert(strSQL string, args ...interface{}) (int64, error) {
	if perr := msd.db.Ping(); perr != nil {
		return 0, perr
	}

	r, err := msd.db.Exec(strSQL, args)
	if err != nil {
		return 0, err
	}

	return r.LastInsertId()
}

// Update : 更新语句
func (msd *MySQLDB) Update(strSQL string, args ...interface{}) (int64, error) {
	if perr := msd.db.Ping(); perr != nil {
		return 0, perr
	}

	r, err := msd.db.Exec(strSQL, args)
	if err != nil {
		return 0, err
	}

	return r.RowsAffected()
}

// Close : 关闭
func (msd *MySQLDB) Close() {
	msd.db.Close()
}
