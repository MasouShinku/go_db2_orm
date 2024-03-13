// 封装事务接口
package session

import "go_db2_orm/db2orm/log"

// Begin 申请事务
func (s *Session) Begin() (err error) {
	log.Infoln("transaction begin")
	if s.tx, err = s.db.Begin(); err != nil {
		log.Errorln(err)
	}
	return
}

// Commit 提交事务
func (s *Session) Commit() (err error) {
	log.Infoln("transaction commit")
	if err = s.tx.Commit(); err != nil {
		log.Errorln(err)
	}
	return
}

// Rollback 回滚事务
func (s *Session) Rollback() (err error) {
	log.Infoln("transaction rollbackS")
	if err = s.tx.Rollback(); err != nil {
		log.Errorln(err)
	}
	return
}
