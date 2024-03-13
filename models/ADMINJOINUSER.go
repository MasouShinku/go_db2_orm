package models

type ADMINJOINUSER struct {
	Id       int `db2orm:"PRIMARY KEY NOT NULL"`
	Name     string
	Desc     string
	UserName string
	Age      int
}
