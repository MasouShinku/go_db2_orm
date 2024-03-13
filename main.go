package main

import (
	_ "github.com/ibmdb/go_ibm_db"
	"go_db2_orm/db2orm"
	"go_db2_orm/db2orm/log"
	"go_db2_orm/models"
)

var con = "HOSTNAME=localhost;DATABASE=mydb;PORT=50000;UID=db2inst1;PWD=ky20021120"

type USER struct {
	Name string `db2orm:"PRIMARY KEY NOT NULL"`
	Age  int
}

func main() {

	engine, _ := db2orm.NewEngine("go_ibm_db", con)
	s := engine.NewSession()
	s.Model(&models.ADMIN{})
	err := s.CreateTable()
	if err != nil {
		log.Errorln(err)
	}

	var user1 = &models.ADMIN{Name: "shinku", Id: 1, Desc: "你是？"}
	var user2 []models.ADMIN

	s.Insert(user1)

	//s.Model(&models.ADMIN{})
	s.Join(&models.ADMIN{}).Limit(2).Where("name = ?", "shinku").Find(&user2)

	log.Infoln(len(user2))

	//s.DropTable()

}
