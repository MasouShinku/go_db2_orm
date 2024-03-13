// 操作数据库表格相关的代码
package session

import (
	"fmt"
	"go_db2_orm/db2orm/log"
	"go_db2_orm/db2orm/schema"
	"reflect"
	"strings"
)

// Model 用于给refTable赋值
func (s *Session) Model(value interface{}) *Session {
	// 已经有refTable或者传入类型不符时不操作
	if s.refTable != nil && reflect.TypeOf(value) == reflect.TypeOf(s.refTable.Model) {
		return s
	}
	// 否则进行赋值操作
	s.refTable = schema.Parse(value, s.dialect)
	return s
}

// RefTable 获取refTable
func (s *Session) RefTable() *schema.Schema {
	var _schema *schema.Schema
	if _schema = s.refTable; _schema == nil {
		log.Errorln("表格未定义!")
	}
	return _schema
}

// 创建表
func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
		log.Infoln(fmt.Sprintf("( %s %s %s ) added to columns", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
	return err
}

// 删除表
func (s *Session) DropTable() error {
	table := s.RefTable()
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", table.Name)).Exec()
	return err
}

// 查询表是否存在
func (s *Session) HasTable() bool {
	table := s.RefTable()
	sql, values := s.dialect.TableExistSQL(table.Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return strings.ToUpper(tmp) == strings.ToUpper(s.RefTable().Name)
}
