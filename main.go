package main

import (
	_ "github.com/ibmdb/go_ibm_db"
	"go_db2_orm/db2orm"
	"go_db2_orm/db2orm/log"
	"go_db2_orm/models"
)

// var con = "HOSTNAME=localhost;DATABASE=mydb;PORT=50000;UID=db2inst1;PWD=ky20021120"
var con = "HOSTNAME=1.94.43.50;DATABASE=testdb;PORT=60000;UID=db2inst1;PWD=db2inst1"

//type USER struct {
//	Name string `db2orm:"PRIMARY KEY NOT NULL"`
//	Age  int
//}

func main() {

	engine, _ := db2orm.NewEngine("go_ibm_db", con)
	s := engine.NewSession()
	s.Model(&models.USER{})
	err := s.CreateTable()
	if err != nil {
		log.Errorln(err)
	}

	var user1 = models.USER{
		Id:       0,
		UserName: "testuser",
		Age:      0,
	}

	//var user1 = &models.ADMIN{Name: "shinku", Id: 1, Desc: "你是？"}
	//var user2 []models.ADMIN

	//var adminjoinuser []models.ADMINJOINUSER

	s.Insert(user1)

	//s.Model(&models.ADMIN{})
	//s.Join(USER{}, "left join").Limit(2).Find(&adminjoinuser)
	//s.Limit(2).Find(&user2)

	//log.Infoln(adminjoinuser)

	//s.DropTable()

}
