package models

// 表格名大写
// 字段名首字母大写

type ADMIN struct {
	Id   int `db2orm:"PRIMARY KEY NOT NULL"`
	Name string
	Desc string
}
