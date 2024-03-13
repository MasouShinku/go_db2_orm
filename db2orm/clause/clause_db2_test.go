package clause

import (
	"reflect"
	"testing"
)

func testSelect_db2(t *testing.T) {
	var clause Clause
	clause.Set(LIMIT, 3)
	clause.Set(SELECT, "USER", []string{"*"})
	clause.Set(WHERE, "Name = ?", "Tom")
	clause.Set(ORDERBY, "Age ASC")
	sql, vars := clause.Build(SELECT, WHERE, ORDERBY, LIMIT)
	t.Log(sql, vars)
	if sql != "SELECT * FROM USER WHERE Name = ? ORDER BY Age ASC LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		t.Fatal("failed to build SQLVars")
	}
}

func TestClause_Build_db2(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		testSelect_db2(t)
	})
}
