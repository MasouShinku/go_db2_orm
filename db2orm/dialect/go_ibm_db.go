package dialect

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type go_ibm_db struct{}

var _ Dialect = (*go_ibm_db)(nil)

func init() {
	RegisterDialect("go_ibm_db", &go_ibm_db{})
}

func (d *go_ibm_db) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.String:
		return "VARCHAR(255)" // DB2 使用 VARCHAR 存储字符串，这里假设一个默认长度，实际应用中可能需要根据实际需求调整
	case reflect.Array, reflect.Slice:
		return "BLOB" // 对应二进制大对象
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return "INTEGER"
	case reflect.Int64, reflect.Uint64:
		return "BIGINT" // 对应较大的整数值
	case reflect.Float32:
		return "REAL" // 单精度浮点数
	case reflect.Float64:
		return "DOUBLE" // 双精度浮点数
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "TIMESTAMP" // DB2中使用TIMESTAMP类型存储日期和时间
		}
	}
	panic(fmt.Sprintf("invalid type %s (%s)", typ.Type().Name(), typ.Kind()))
}

func (d *go_ibm_db) TableExistSQL(tableName string) (string, []interface{}) {
	// db2一定要转大写
	tableName = strings.ToUpper(tableName)
	args := []interface{}{tableName}
	sql := "SELECT tabname FROM syscat.tables WHERE tabschema = CURRENT SCHEMA AND tabname = ?"
	return sql, args
}
