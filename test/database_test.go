package test

import (
	"http_server/db"
	"os"
	"testing"
)

func TestDatabase(t *testing.T) {

	os.Setenv("database_path", "127.0.0.1:8899")

	dbsesssion := db.DBsession()
	userID := dbsesssion.InsertUser("sample_name", "sample_password")
	alluser := dbsesssion.Find()

	defer dbsesssion.RemoveAll()

	if len(alluser) != 1 {

		t.Fatalf("not inserted")
	}

	if alluser[0].ID != userID {

		t.Fatalf("ID is not matched")
	}

}
