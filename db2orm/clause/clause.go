// 用于拼接独立子句
package clause

import (
	"strings"
)

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
	JOIN
)

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

// Set 获取操作类型对应的参数
func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
	//log.Infoln(fmt.Sprintf("there is vars after Set : [%s] %v", c.sql[name], c.sqlVars[name]))
}

// Build 生成sql语句
func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		//log.Infoln(fmt.Sprintf("building %s ...", order))
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
		//log.Infoln(fmt.Sprintf("built sql :  %s", sqls))
	}
	return strings.Join(sqls, " "), vars
}
