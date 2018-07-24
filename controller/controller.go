package controller

import (
	"encoding/json"
	"fmt"
	"http_server/db"
	"http_server/middleware"
	"http_server/model"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

var session = db.DBsession()

func FriendlistHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("friendlist")
	userid := middleware.RequestIDFromContext(r.Context())
	print(userid)
	myuser := session.FindOneById(userid)
	friendIds := myuser.FriendIds
	friendsdata := session.FindById(friendIds)

	FriendStruct := []model.Friend{}

	for i := 0; i < len(friendsdata); i++ {
		FriendStruct = append(FriendStruct,
			model.Friend{Name: friendsdata[i].Name,
				ID:            friendsdata[i].ID.Hex(),
				Photourl:      friendsdata[i].Photourl,
				Backgroundurl: friendsdata[i].Backgroundurl})
	}
	json.NewEncoder(w).Encode(FriendStruct)
}

func FindFriend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	friendname := vars["friendname"]
	friend := session.FindOneByName(friendname)

	json.NewEncoder(w).Encode(model.Friend{
		Name:          friend.Name,
		ID:            friend.ID.Hex(),
		Photourl:      friend.Photourl,
		Backgroundurl: friend.Backgroundurl},
	)
}

func AddFriend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	friendid := vars["friendid"]
	userid := middleware.RequestIDFromContext(r.Context())

	session.AddFriend(userid, friendid)
	session.AddFriend(friendid, userid)
}

func GetProfile(w http.ResponseWriter, r *http.Request) {

	userid := middleware.RequestIDFromContext(r.Context())
	user := session.FindOneById(userid)
	profile := model.Profile{Name: user.Name, Photourl: user.Photourl}
	json.NewEncoder(w).Encode(profile)

}

func ProfilePhoto(w http.ResponseWriter, r *http.Request) {
	userid := middleware.RequestIDFromContext(r.Context())

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("image")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	filename := uuid.Must(uuid.NewV4()).String()

	f, err := os.OpenFile("./static/img/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	session.UpdateOne("photourl", "http://localhost:8181/static/img/"+filename, userid)
}

func ProfileName(w http.ResponseWriter, r *http.Request) {
	userid := middleware.RequestIDFromContext(r.Context())

	decoder := json.NewDecoder(r.Body)
	var profile model.Profile
	err := decoder.Decode(&profile)

	if err != nil {
		panic(err)
	}

	session.UpdateOne("name", profile.Name, userid)
}

func Storedtalks(w http.ResponseWriter, r *http.Request) {
	userid := middleware.RequestIDFromContext(r.Context())

	user := session.FindOneById(userid)
	json.NewEncoder(w).Encode(user.Talks)
	session.UpdateOne("talks", []model.Talk{}, userid)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var loginform model.LoginForm
		err := decoder.Decode(&loginform)

		if err != nil {
			panic(err)
		}

		name := loginform.Name
		password := loginform.Password

		user := session.GetByName(name)
		fmt.Printf(user.ID.Hex())

		if user == nil {
			userid := session.InsertUser(name, password)

			// Set some claims
			token := jwt.New(jwt.SigningMethodRS512)
			claims := make(jwt.MapClaims)
			claims["exp"] = time.Now().Add(time.Minute * 50000).Unix()
			claims["id"] = userid
			token.Claims = claims
			tokenString, err := token.SignedString([]byte("secret"))

			if err != nil {

			}

			login := model.Login{Login: true, Token: tokenString, ID: userid.Hex(), Name: name}
			json.NewEncoder(w).Encode(login)

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
				json.NewEncoder(w).Encode(login)

			} else {
				login := model.Login{Login: false}
				json.NewEncoder(w).Encode(login)
			}
		}

	}
}

func FeedController(w http.ResponseWriter, r *http.Request) {
	userid := middleware.RequestIDFromContext(r.Context())

	if r.Method == "POST" {

		// userid := middleware.RequestIDFromContext(r.Context())

		r.ParseMultipartForm(32 << 20)
		file, handler, _ := r.FormFile("image")

		fmt.Fprintf(w, "%v", handler.Header)

		defer file.Close()

		feedcontent := r.Form["feedcontent"][0]
		print(feedcontent)

		filename := uuid.Must(uuid.NewV4()).String()
		photourl := "http://localhost:8181/static/img/" + filename

		f, err := os.OpenFile("./static/img/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		session.InsertFeed(feedcontent, userid, photourl)
	} else {
		user := session.FindOneById(userid)
		friendids := []string{}
		for i := 0; i < len(user.FriendIds); i++ {
			friendids = append(friendids, user.FriendIds[i].Hex())
		}
		feeds := session.GetFeed(friendids)

		json.NewEncoder(w).Encode(feeds)

	}
}
