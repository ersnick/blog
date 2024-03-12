package transport

import (
	handlers "blog/handlers"
	"blog/handlers/authorization"
	"blog/handlers/likes"
	"blog/handlers/post"
	"blog/handlers/profile"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	StaticDir = "/images/"
)

func RoutingHandler() {
	router := mux.NewRouter()
	router.HandleFunc("/", handlers.IndexPageHandler).Methods("GET")
	router.HandleFunc("/", handlers.Indexhandler).Methods("POST")

	router.HandleFunc("/login", authorization.LoginPageHandler).Methods("GET")
	router.HandleFunc("/login", authorization.LoginHandler).Methods("POST")
	router.HandleFunc("/logout", authorization.Logout)
	router.HandleFunc("/delete", authorization.DeleteHandler)

	router.HandleFunc("/signup", authorization.SignUpPageHandler).Methods("GET")
	router.HandleFunc("/signup", authorization.SignUpHandler).Methods("POST")

	router.HandleFunc("/update", handlers.UpdateHandler).Methods("POST")
	router.HandleFunc("/update", handlers.UpdatePage).Methods("GET")
	router.HandleFunc("/updatePost", post.UpdatePost).Methods("POST")
	router.HandleFunc("/updatePost", post.UpdatePostPage).Methods("GET")
	router.HandleFunc("/deletePost", post.DeletePostHandler)
	router.HandleFunc("/users", handlers.UsersPageHandler)

	router.HandleFunc("/profile", profile.ProfilePage).Methods("GET")
	router.HandleFunc("/profile", profile.ProfilePageInputHandler).Methods("POST")

	router.HandleFunc("/like", likes.LikeHandler)
	router.HandleFunc("/unlike", likes.UnlikeHandler)

	router.
		PathPrefix(StaticDir).
		Handler(http.StripPrefix(StaticDir, http.FileServer(http.Dir("."+StaticDir))))
	log.Fatalln(http.ListenAndServe("localhost:8000", router))
}
