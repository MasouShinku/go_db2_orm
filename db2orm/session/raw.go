/*
	原生交互部分
*/

package session

import (
	"database/sql"
	"go_db2_orm/db2orm/clause"
	"go_db2_orm/db2orm/dialect"
	"go_db2_orm/db2orm/log"
	"go_db2_orm/db2orm/schema"
	"strings"
)

// Session 会话结构体
type Session struct {
	db       *sql.DB         //数据库指针
	sql      strings.Builder // 用于存放待执行的sql语句
	sqlVars  []interface{}   // 用于存储拼接sql语句时的占位符
	dialect  dialect.Dialect // 存放数据库方言
	refTable *schema.Schema  // 数据库模式
	clause   clause.Clause   // sql子句
	tx       *sql.Tx         // 事务锁
}

// New 创建会话对象
func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

// Clear 清空sql对象
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

// GetDB 获取DB对象
func (s *Session) GetDB() *sql.DB {
	return s.db
}

// Raw 创建sql语句
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// Exec sql的执行
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Infoln(s.sql.String(), s.sqlVars)
	if result, err = s.GetDB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Errorln(err)
	}
	return
}

// QueryRow 获取单条结果
func (s *Session) QueryRow() (row *sql.Row) {
	defer s.Clear()
	log.Infoln(s.sql.String(), s.sqlVars)
	row = s.GetDB().QueryRow(s.sql.String(), s.sqlVars...)
	return
}

// QueryRows 获取多条结果
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Infoln(s.sql.String(), s.sqlVars)
	if rows, err = s.GetDB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Errorln(err)
	}
	return
}

// DB的func集合
type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}
