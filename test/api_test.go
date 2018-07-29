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

func TestLoginHandler(t *testing.T) {
	url := "http://localhost:8181/api/login/"

	dbsession := db.DBsession()
	dbsession.RemoveAll()
	defer dbsession.RemoveAll()
	fmt.Println("URL:>", url)

	loginmodel := &model.LoginForm{Name: "newuser", Password: "newpass"}

	loginjson, err := json.Marshal(loginmodel)
	if err != nil {
		fmt.Println(err)
		return
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

	if err != nil {
		panic(err)
	}

	if loginreturn.Login != true {
		t.Fatalf("ID is not matched")
	}

	users := dbsession.Find()

	if len(users) != 1 {
		t.Fatalf("number of user is not correct")
	}
}
