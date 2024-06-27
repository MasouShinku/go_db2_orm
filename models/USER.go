package models

// 表格名大写
// 字段名首字母大写

type USER struct {
	Id       int `db2orm:"PRIMARY KEY NOT NULL"`
	UserName string
	Age      int
}
