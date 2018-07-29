package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"http_server/db"
	"http_server/model"
	"io/ioutil"
	"net/http"
	"testing"
)

func loginuser(name string, password string) *model.Login {
	url := "http://localhost:8181/api/login/"

	loginmodel := &model.LoginForm{Name: name, Password: password}

	loginjson, err := json.Marshal(loginmodel)
	if err != nil {
		fmt.Println(err)

	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(loginjson))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	loginreturn := &model.Login{}
	err = json.Unmarshal(body, loginreturn)

	return loginreturn
}

func TestLoginHandler(t *testing.T) {
	url := "http://localhost:8181/api/login/"

	dbsession := db.DBsession()
	dbsession.RemoveAll()
	defer dbsession.RemoveAll()
	fmt.Println("URL:>", url)

	loginreturn := loginuser("newuser", "password")

	if loginreturn.Login != true {
		t.Fatalf("ID is not matched")
	}

	users := dbsession.Find()

	if len(users) != 1 {
		t.Fatalf("number of user is not correct")
	}

	loginreturn2 := loginuser("newuser2", "password2")

	if loginreturn2.Login != true {
		t.Fatalf("ID is not matched")
	}

	users2 := dbsession.Find()

	if len(users2) != 2 {
		t.Fatalf("number of user is not correct")
	}

	loginreturn3 := loginuser("newuser", "password")

	if loginreturn3.Login != true {
		t.Fatalf("ID is not matched")
	}

	users3 := dbsession.Find()

	if len(users3) != 2 {
		t.Fatalf("number of user is not correct")
	}

	if loginreturn3.ID != loginreturn.ID {
		t.Fatalf("ID is not matched")
	}

	loginreturn4 := loginuser("newuser", "password_wrong")

	if loginreturn4.Login != false {
		t.Fatalf("password ")
	}

	users4 := dbsession.Find()

	if len(users4) != 2 {
		t.Fatalf("number of user is not correct")
	}

}
