package main

import (
	"http_server/controller"
	"http_server/db"
	"http_server/message"
	"http_server/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// fs := http.FileServer(http.Dir("static"))
	db.DBsession()
	r := mux.NewRouter()
	// r.Handle("/", fs)
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/api/login/", controller.LoginHandler)
	r.Handle("/api/friendslist/", middleware.Middleware(controller.FriendlistHandler))
	r.HandleFunc("/websocket", message.Websocketmessage)
	r.HandleFunc("/api/find/{friendname}", controller.FindFriend)
	r.Handle("/api/addfriend/{friendid}", middleware.Middleware(controller.AddFriend))
	r.Handle("/api/profile", middleware.Middleware(controller.GetProfile))
	r.Handle("/api/profile/photo", middleware.Middleware(controller.ProfilePhoto))
	r.Handle("/api/feed", middleware.Middleware(controller.FeedController))
	r.Handle("/api/storedtalks", middleware.Middleware(controller.Storedtalks))

	err := http.ListenAndServe(":8181", r)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
