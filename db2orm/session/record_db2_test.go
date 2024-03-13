package session

//
//import "testing"
//
//type USER struct {
//	Name string `db2orm:"PRIMARY KEY NOT NULL"`
//	Age  int
//}
//
//var (
//	user1 = &USER{"Tom", 18}
//	user2 = &USER{"Sam", 25}
//	user3 = &USER{"Jack", 25}
//)
//
//func testRecordInit_db2(t *testing.T) *Session {
//	t.Helper()
//	s := NewSession().Model(&USER{})
//	err1 := s.DropTable()
//	err2 := s.CreateTable()
//	_, err3 := s.Insert(user1, user2)
//	if err1 != nil {
//		t.Fatal("err1 occurred")
//	}
//	if err2 != nil {
//		t.Fatal("err2 occurred")
//	}
//	if err3 != nil {
//		t.Fatal("err3 occurred")
//	}
//	if err1 != nil || err2 != nil || err3 != nil {
//		t.Fatal("failed init test records")
//	}
//	return s
//}
//
//func TestSession_Insert_db2(t *testing.T) {
//	s := testRecordInit_db2(t)
//	affected, err := s.Insert(user3)
//	if err != nil || affected != 1 {
//		t.Fatal("failed to create record")
//	}
//}
//
//func TestSession_Find_db2(t *testing.T) {
//	s := testRecordInit_db2(t)
//	var users []USER
//	if err := s.Find(&users); err != nil || len(users) != 2 {
//		t.Fatal("failed to query all")
//	}
//}
//
//func TestSession_Limit_db2(t *testing.T) {
//	s := testRecordInit_db2(t)
//	var users []USER
//	err := s.Limit(1).Find(&users)
//	if err != nil || len(users) != 1 {
//		t.Fatal("failed to query with limit condition")
//	}
//}
//
//func TestSession_Update_db2(t *testing.T) {
//	s := testRecordInit_db2(t)
//	affected, _ := s.Where("Name = ?", "Tom").Update("Age", 30)
//	u := &USER{}
//	_ = s.OrderBy("Age DESC").First(u)
//
//	if affected != 1 || u.Age != 30 {
//		t.Fatal("failed to update")
//	}
//}
//
//func TestSession_DeleteAndCount_db2(t *testing.T) {
//	s := testRecordInit_db2(t)
//	affected, _ := s.Where("Name = ?", "Tom").Delete()
//	count, _ := s.Count()
//
//	if affected != 1 || count != 1 {
//		t.Fatal("failed to delete or count")
//	}
//}
