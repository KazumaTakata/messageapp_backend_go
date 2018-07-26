package db

import (
	"fmt"
	"http_server/model"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type dbsession struct {
	session *mgo.Session
}

func (db *dbsession) AddFriend(myid string, friendid string) {

	c := db.session.DB("swiftline").C("users")

	query := bson.M{"_id": bson.ObjectIdHex(myid)}
	update := bson.M{"$push": bson.M{"friendIds": bson.ObjectIdHex(friendid)}}

	// Update
	err := c.Update(query, update)
	if err != nil {
		fmt.Println(err)
	}
}

func (db *dbsession) UpdateOne(fieldname string, fieldvalue interface{}, myid string) {

	c := db.session.DB("swiftline").C("users")

	query := bson.M{"_id": bson.ObjectIdHex(myid)}
	update := bson.M{"$set": bson.M{fieldname: fieldvalue}}

	// Update
	err := c.Update(query, update)
	if err != nil {
		fmt.Println(err)
	}
}

func (db *dbsession) FindOneByName(name string) model.User {

	c := db.session.DB("swiftline").C("users")

	var result model.User
	err := c.Find(bson.M{"name": name}).One(&result)

	if err != nil {
		fmt.Println(err)
	}

	return result
}

func (db *dbsession) FindOneById(userid string) model.User {

	c := db.session.DB("swiftline").C("users")

	var result model.User
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(userid)}).One(&result)

	if err != nil {
		fmt.Println(err)
	}

	return result
}

func (db *dbsession) FindById(friendIds []bson.ObjectId) []model.User {
	c := db.session.DB("swiftline").C("users")

	var result []model.User
	err := c.Find(bson.M{"_id": bson.M{"$in": friendIds}}).All(&result)
	if err != nil {

	}

	return result

}

func (db *dbsession) GetByName(name string) *model.User {
	c := db.session.DB("swiftline").C("users")

	var result model.User
	err := c.Find(bson.M{"name": name}).One(&result)
	if err != nil {
		return nil
	}
	return &result
}

func (db *dbsession) InsertUser(name string, password string) bson.ObjectId {
	c := db.session.DB("swiftline").C("users")
	i := bson.NewObjectId()

	err := c.Insert(&model.User{ID: i, Name: name, Password: password,
		Photourl:      "http://localhost:8181/img/defaultprofile.png",
		Backgroundurl: "http://localhost:8181/img/rocco-caruso-722282-unsplash.jpg",
		Talks:         []model.Talk{},
		FriendIds:     []bson.ObjectId{},
	})

	if err != nil {
		fmt.Println(err)
	}

	return i
}

func (db *dbsession) InsertFeed(feedcontent string, userid string, photourl string) {
	c := db.session.DB("swiftline").C("feeds")

	err := c.Insert(&model.Feed{Feedcontent: feedcontent, UserId: userid, Photourl: photourl})

	if err != nil {
		fmt.Println(err)
	}
}

func (db *dbsession) GetFeed(friendids []string) []model.Feed {
	c := db.session.DB("swiftline").C("feeds")

	var result []model.Feed
	err := c.Find(bson.M{"userId": bson.M{"$in": friendids}}).All(&result)
	if err != nil {

	}

	if err != nil {
		fmt.Println(err)
	}

	return result
}

func (db *dbsession) InsertTalk(friendid string, content string, myId string) {
	c := db.session.DB("swiftline").C("users")

	query := bson.M{"_id": bson.ObjectIdHex(friendid)}
	update := bson.M{"$push": bson.M{"talks": bson.M{"content": content, "id": myId}}}

	// Update
	err := c.Update(query, update)
	if err != nil {
		fmt.Println(err)
	}

}

var (
	sessionIns *dbsession
)

func DBsession() *dbsession {

	if sessionIns == nil {
		// dbaddress := "127.0.0.1:27017"
		dbaddress := "messenger_mongo:27017"

		session, err := mgo.Dial(dbaddress)
		if err != nil {
			fmt.Println(err)
		}
		sessionIns = &dbsession{
			session: session,
		}
	}
	return sessionIns
}
