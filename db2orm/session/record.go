// 增删改查相关代码
package session

import (
	"errors"
	"fmt"
	"go_db2_orm/db2orm/clause"
	"go_db2_orm/db2orm/log"
	"go_db2_orm/db2orm/schema"
	"reflect"
)

// Insert 期望操作形式：
// u1:=&User{Name:"shinku" ,Age:18}
// u2:=&User{Name:"shizuka",Age:20}
// s.Insert(u1,u2)
// 实际sql语句:
// INSERT INTO $tableName($col1,$col2,...) VALUES
// (A1,A2,...),
// (B1,B2,...),
// ...
func (s *Session) Insert(values ...interface{}) (int64, error) {

	// 设置INSERT部分
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		s.CallMethod(BeforeInsert, value)
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}

	// 设置VALUES部分
	s.clause.Set(clause.VALUES, recordValues...)

	// 拼接生成完整语句
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterInsert, nil)
	return result.RowsAffected()
}

// Find 查询语句
// 期望调用形式：
// var users []User
// s.Find(&users)
func (s *Session) Find(values interface{}) error {
	s.CallMethod(BeforeQuery, nil)
	destSlice := reflect.Indirect(reflect.ValueOf(values))

	// 获取元素类型，并映射出表结构
	destType := destSlice.Type().Elem()
	//table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()
	table := s.RefTable()
	// 拼接出select语句并执行
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames, s.clause.JoinList)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	//sql, vars := s.clause.Build(clause.SELECT)
	//log.Infoln(fmt.Sprintf("here is sql : %s", sql))
	rows, err := s.Raw(sql, vars...).QueryRows()

	//columns, err := rows.Columns()
	//for _, column := range columns {
	//	fmt.Printf(column)
	//}

	if err != nil {
		return err
	}

	tempRefTable := schema.Parse(reflect.New(destType).Interface(), s.dialect)

	// 将结果平铺，遍历后按字段赋值
	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range tempRefTable.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}

		if err := rows.Scan(values...); err != nil {
			log.Errorln(err)
			return err
		}
		s.CallMethod(AfterQuery, dest.Addr().Interface())
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	log.Infoln(destSlice.Len())
	return rows.Close()
}

// Update 更新语句
// 支持入参1：map[string]any的键值对
// 支持入参2：kv列表：key1,val1,kay2,val2,...
func (s *Session) Update(kv ...interface{}) (int64, error) {
	s.CallMethod(BeforeUpdate, nil)
	m, ok := kv[0].(map[string]interface{})
	// 若不是map形式，则对平铺kv列表进行转换
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}

	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterUpdate, nil)
	return result.RowsAffected()
}

// Delete 删除语句
func (s *Session) Delete() (int64, error) {
	s.CallMethod(BeforeDelete, nil)
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterDelete, nil)
	return result.RowsAffected()
}

// Count 计数语句
func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

// Limit 限制结果数量
// 返回Session指针，方便链式调用
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

// Where 设置条件
// 返回Session指针，方便链式调用
func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

// OrderBy 设置排序
// 返回Session指针，方便链式调用
func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

// First 返回第一条记录
func (s *Session) First(value interface{}) error {
	desc := reflect.Indirect(reflect.ValueOf(value))
	descSlice := reflect.New(reflect.SliceOf(desc.Type())).Elem()
	if err := s.Limit(1).Find(descSlice.Addr().Interface()); err != nil {
		return err
	}
	if descSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	desc.Set(descSlice.Index(0))
	return nil
}

// Join 第一个参数是struct,第二个参数是joinType,第三个参数是cond
func (s *Session) Join(value interface{}, joinType string, cond ...string) *Session {
	valueType := reflect.TypeOf(value)
	if valueType.Kind() == reflect.Ptr {
		valueType = valueType.Elem()
	}
	tabname := valueType.Name()
	//s.clause.Set(clause.JOIN, tabname)

	// 更新join数组
	var newJoinInfo clause.JoinInfo
	newJoinInfo.JoinType = joinType
	if len(cond) > 0 {
		newJoinInfo.OnCond = cond[0]
	}
	newJoinInfo.Tablename = tabname
	for i := 0; i < valueType.NumField(); i++ {
		fieldName := valueType.Field(i).Name
		newJoinInfo.FormattedVals = append(newJoinInfo.FormattedVals, fmt.Sprintf("%s.%s", tabname, fieldName))
	}

	s.clause.JoinList = append(s.clause.JoinList, newJoinInfo)

	//log.Infoln(fmt.Sprintf("%v", newJoinInfo))

	// 传送点
	// 这里需要重新设置一下select的参数

	return s
}
