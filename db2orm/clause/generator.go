// 子句的生成规则
package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderby
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
	generators[JOIN] = _join
}

func _insert(values ...interface{}) (string, []interface{}) {
	// INSERT INTO $tablename ($fields)
	tablename := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tablename, fields), []interface{}{}
}

func _values(values ...interface{}) (string, []interface{}) {
	// VALUES ($v1),($v2),...
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			var tempVars []string
			for j := 0; j < len(v); j++ {
				tempVars = append(tempVars, "?")
			}
			bindStr = strings.Join(tempVars, ", ")
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

func _select(values ...interface{}) (string, []interface{}) {
	// SELECT $fields FROM $tablename
	tablename := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tablename), []interface{}{}
}

func _limit(values ...interface{}) (string, []interface{}) {
	// LIMIT $num
	return fmt.Sprintf("LIMIT ?"), values
}

func _where(values ...interface{}) (string, []interface{}) {
	// WHERE $desc
	desc := values[0]
	vars := values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func _orderby(values ...interface{}) (string, []interface{}) {
	// ORDER BY $field
	field := values[0]
	return fmt.Sprintf("ORDER BY %s", field), []interface{}{}
}

func _update(values ...interface{}) (string, []interface{}) {
	tablename := values[0]
	m := values[1].(map[string]interface{})
	var keys []string
	var vars []interface{}
	for k, v := range m {
		keys = append(keys, k+" =? ")
		vars = append(vars, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tablename, strings.Join(keys, ", ")), vars
}

func _delete(values ...interface{}) (string, []interface{}) {
	tablename := values[0]
	return fmt.Sprintf("DELETE FROM %s", tablename), []interface{}{}
}

func _count(values ...interface{}) (string, []interface{}) {
	tablename := values[0]
	return _select(tablename, []string{"count(*)"})
}

func _join(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("JOIN %s", values[0]), []interface{}{}
}
