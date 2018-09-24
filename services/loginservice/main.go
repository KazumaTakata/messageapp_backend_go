package main

import (
	"context"
	"fmt"
	"log"

	"http_server/db"
	"http_server/model"
	pb "http_server/services/loginservice/proto"

	jwt "github.com/dgrijalva/jwt-go"
	micro "github.com/micro/go-micro"
)

var session = db.DBsession()

type ILogin interface {
	DoLogin(*pb.Userdata) (*model.Login, error)
}

type Login struct {
}

func (repo *Login) DoLogin(userdata *pb.Userdata) (*model.Login, error) {

	name := userdata.Name
	password := userdata.Password

	user := session.FindOneByName(name)
	// fmt.Printf(user.ID.Hex())

	if user == nil {
		userid := session.InsertUser(name, password)

		// Set some claims
		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &model.Token{
			ID: userid.Hex(),
		})

		tokenString, err := token.SignedString([]byte("secret"))

		if err != nil {

		}

		login := model.Login{Login: true, Token: tokenString, ID: userid.Hex(), Name: name}
		return &login, nil

	} else {
		if user.Password == password {

			token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &model.Token{
				ID: user.ID.Hex(),
			})

			tokenstring, err := token.SignedString([]byte("secret"))
			if err != nil {
				log.Fatalln(err)
			}

			if err != nil {

			}

			login := model.Login{Login: true, Token: tokenstring, ID: user.ID.Hex(), Name: name}
			return &login, nil

		} else {
			login := model.Login{Login: false}
			return &login, nil
		}
	}

}

type service struct {
	login ILogin
}

func (s *service) LoginOrSignup(ctx context.Context, req *pb.Userdata, res *pb.Response) error {
	login, err := s.login.DoLogin(req)

	if err != nil {
		return err
	}

	res.Id = login.ID
	res.Name = login.Name
	res.Login = login.Login
	res.Token = login.Token
	return nil
}

func main() {
	login := &Login{}

	srv := micro.NewService(

		// This name must match the package name given in your protobuf definition
		micro.Name("go.micro.srv.login"),
		micro.Version("latest"),
	)

	srv.Init()

	pb.RegisterLoginServiceHandler(srv.Server(), &service{login})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}

}
