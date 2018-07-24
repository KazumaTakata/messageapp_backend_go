package model

import (
	jwt "github.com/dgrijalva/jwt-go"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Friend struct {
	Name          string `json:"name"`
	Photourl      string `json:"photourl"`
	ID            string `json:"id"`
	Backgroundurl string `json:"backgroundurl"`
}

type Feed struct {
	Feedcontent string `json:"feedcontent" bson:"feedcontent"`
	UserId      string `json:"userId" bson:"userId"`
	Photourl    string `json:"photourl" bson:"photourl"`
}

type Profile struct {
	Name     string `json:"name"`
	Photourl string `json:"photourl"`
}

type Token struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

type LoginForm struct {
	Name     string
	Password string
}

type Login struct {
	Login bool   `json:"login"`
	Token string `json:"token"`
	ID    string `json:"id"`
	Name  string `json:"name"`
}

type Talk struct {
	Content string `json:"content"`
	ID      string `json:"id"`
}

type Id struct {
	ID bson.ObjectId `bson:"_id"`
}

type User struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	Name          string
	Password      string
	Photourl      string
	Backgroundurl string
	FriendIds     []bson.ObjectId `bson:"friendIds"`
	Talks         []Talk
}

type dbsession struct {
	session *mgo.Session
}
