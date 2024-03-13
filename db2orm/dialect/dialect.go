// 数据库方言部分,为多种数据库提供类型映射
package dialect

import (
	"reflect"
)

type Dialect interface {
	// 将GO语言的类型转换为数据库的数据类型
	DataTypeOf(typ reflect.Value) string
	// 返回查询表是否存在的SQL
	TableExistSQL(tableName string) (string, []interface{})
}

var dialectsMap = map[string]Dialect{}

// 方言的注册部分
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// 获取方言映射表
func GetDialect(name string) (dialect Dialect, isExist bool) {
	dialect, isExist = dialectsMap[name]
	return
}
