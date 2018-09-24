package main

import (
	"context"
	"encoding/json"
	"fmt"
	proto "http_server/services/loginservice/proto"
	"io/ioutil"

	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
)

func parseFile(file string) (*proto.Userdata, error) {
	var userdata *proto.Userdata
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &userdata)
	return userdata, err
}

func main() {

	cmd.Init()

	// Create a new service. Optionally include some options here.

	// Create new greeter client
	client := proto.NewLoginServiceClient("go.micro.srv.login", microclient.DefaultClient)
	userdata, _ := parseFile("data.json")
	// Call the greeter

	login, err := client.LoginOrSignup(context.TODO(), userdata)

	if err != nil {
	}

	fmt.Printf(login.Name)

}
