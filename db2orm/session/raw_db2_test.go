package session

import (
	"database/sql"
	"fmt"
	"go_db2_orm/db2orm/dialect"
	"go_db2_orm/db2orm/log"
	"os"
	"testing"

	_ "github.com/ibmdb/go_ibm_db"
	_ "github.com/mattn/go-sqlite3"
)

var (
	TestDB      *sql.DB
	TestDial, _ = dialect.GetDialect("go_ibm_db")
)

func TestMain(m *testing.M) {
	con := "HOSTNAME=localhost;DATABASE=mydb;PORT=50000;UID=db2inst;PWD=ky20021120"
	TestDB, err := sql.Open("go_ibm_db", con)

	if err != nil {
		log.Errorln(fmt.Sprintf("connect fail : %s", err))
	}

	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewSession() *Session {
	return New(TestDB, TestDial)
}

func TestSession_Exec_db2(t *testing.T) {
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS USER;").Exec()
	_, _ = s.Raw("CREATE TABLE USER(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	if count, err := result.RowsAffected(); err != nil || count != 2 {
		t.Fatal("expect 2, but got", count)
	}
}

func TestSession_QueryRows_db2(t *testing.T) {
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	row := s.Raw("SELECT count(*) FROM User").QueryRow()
	var count int
	if err := row.Scan(&count); err != nil || count != 0 {
		t.Fatal("failed to query db", err)
	}
}
