// Engine模块
// 交互前的准备工作以及交互后的收尾工作
package db2orm

import (
	"database/sql"
	"go_db2_orm/db2orm/dialect"
	"go_db2_orm/db2orm/log"
	"go_db2_orm/db2orm/session"
)

// Engine 结构体定义
type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

// NewEngine 创建实例方法
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Errorln(err)
		return
	}
	if err = db.Ping(); err != nil {
		log.Errorln(err)
		return
	}

	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorln("dialect %s 不存在!", driver)
		return
	}
	e = &Engine{db: db, dialect: dial}
	log.Infoln("数据库连接成功!")
	return
}

// Close 关闭数据库连接
func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Errorln("数据库关闭失败!")
	}
	log.Infoln("数据库关闭成功!")
}

// NewSession 创建会话
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

// 定义事务方法
type TxFunc func(*session.Session) (any, error)

func (e *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			s.Rollback()
			panic(p)
		} else if err != nil {
			s.Rollback()
		} else {
			err = s.Commit()
		}
	}()

	return f(s)
}
