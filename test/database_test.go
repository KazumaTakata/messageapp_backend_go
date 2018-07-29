package test

import (
	"http_server/config"
	"http_server/db"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestInsertUser(t *testing.T) {
	config.Database_path = "127.0.0.1:8899"
	dbsession := db.DBsession()
	defer dbsession.RemoveAll()

	userID := dbsession.InsertUser("sample_name", "sample_password")
	alluser := dbsession.Find()

	if len(alluser) != 1 {
		t.Fatalf("not inserted")
	}

	if alluser[0].ID != userID {
		t.Fatalf("ID is not matched")
	}

	userID2 := dbsession.InsertUser("sample_name2", "sample_password2")
	alluser2 := dbsession.Find()

	if len(alluser2) != 2 {
		t.Fatalf("not inserted")
	}

	if alluser2[1].ID != userID2 {
		t.Fatalf("ID is not matched")
	}
}

func TestFindOneByID(t *testing.T) {

	config.Database_path = "127.0.0.1:8899"
	dbsession := db.DBsession()
	defer dbsession.RemoveAll()

	userID := dbsession.InsertUser("sample_name", "sample_password")

	user := dbsession.FindOneById(userID.Hex())

	if user.ID != userID {
		t.Fatalf("userid is not matched")
	}

}

func TestFindOneByName(t *testing.T) {

	config.Database_path = "127.0.0.1:8899"
	dbsession := db.DBsession()
	defer dbsession.RemoveAll()

	userID := dbsession.InsertUser("sample_name", "sample_password")

	user := dbsession.FindOneByName("sample_name")

	if user.ID != userID {
		t.Fatalf("userid is not matched")
	}
}

func TestAddFriend(t *testing.T) {

	config.Database_path = "127.0.0.1:8899"
	dbsession := db.DBsession()
	defer dbsession.RemoveAll()

	userID := dbsession.InsertUser("sample_name", "sample_password")
	userID2 := dbsession.InsertUser("sample_name2", "sample_password2")

	dbsession.AddFriend(userID.Hex(), userID2.Hex())

	user := dbsession.FindOneByName("sample_name")
	if user.FriendIds[0] != userID2 {
		t.Fatalf("friend is not added")
	}
}

func TestUpdateOne(t *testing.T) {

	config.Database_path = "127.0.0.1:8899"
	dbsession := db.DBsession()
	defer dbsession.RemoveAll()

	userID := dbsession.InsertUser("sample_name", "sample_password")

	dbsession.UpdateOne("name", "sample_name11", userID.Hex())

	user := dbsession.FindOneById(userID.Hex())

	if user.Name != "sample_name11" {
		t.Fatalf("update is not correct")
	}
}

func TestFindById(t *testing.T) {

	config.Database_path = "127.0.0.1:8899"
	dbsession := db.DBsession()
	defer dbsession.RemoveAll()

	userID := dbsession.InsertUser("sample_name", "sample_password")
	userID2 := dbsession.InsertUser("sample_name2", "sample_password2")

	users := dbsession.FindById([]bson.ObjectId{userID, userID2})

	if len(users) != 2 {
		t.Fatalf("data length is not correct")
	}

	if users[0].ID != userID && users[1].ID != userID2 {
		t.Fatalf("user retrieval is not correct")
	}
}

func TestInsertFeed(t *testing.T) {

	config.Database_path = "127.0.0.1:8899"
	dbsession := db.DBsession()
	defer dbsession.RemoveAll()

	userID := dbsession.InsertUser("sample_name", "sample_password")

	dbsession.InsertFeed("feedcontent1", userID.Hex(), "oeoeoe")

	feeds := dbsession.GetFeed([]string{userID.Hex()})

	if feeds[0].Feedcontent != "feedcontent1" {
		t.Fatalf("insertfeed is not correct")
	}

}

func TestInsertTalk(t *testing.T) {

	config.Database_path = "127.0.0.1:8899"
	dbsession := db.DBsession()
	defer dbsession.RemoveAll()

	userID := dbsession.InsertUser("sample_name", "sample_password")
	userID2 := dbsession.InsertUser("sample_name2", "sample_password2")

	dbsession.InsertTalk(userID.Hex(), "talk1", userID2.Hex())

	user := dbsession.FindOneById(userID.Hex())

	if len(user.Talks) != 1 {
		t.Fatalf("talk data length is not correct")
	}

	if user.Talks[0].Content != "talk1" {
		t.Fatalf("inserted talk content is not correct")
	}

}
