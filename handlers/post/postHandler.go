package post

import (
	"blog/data"
	"blog/db"
	"blog/handlers/authorization"
	"blog/models"
	"html/template"
	"log"
	"net/http"
)

func UpdatePostPage(w http.ResponseWriter, r *http.Request) {
	tplPostUpdate, err := template.ParseFiles("./templates/updatePost.html", "./templates/base.html")
	if err != nil {
		log.Println(err)
	}
	postToBeUpdated := r.URL.Query().Get("id")
	rows, err := db.Db.Query("SELECT * FROM Posts where id=$1", postToBeUpdated)
	if err != nil {
		panic(err)
	}
	var id int64
	var username string
	var body string
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &body, &username)
		if err != nil {
			panic(err)
		}
		log.Println(body)
	}
	var post models.Post
	post.Id = id
	post.Body = body
	post.Username = username

	var data data.UpdatePageData
	data.Post = post
	data.LoggedUser = authorization.Authorized_user
	if len(authorization.Authorized_user) < 1 {
		data.IsLoggedInUser = false
	} else {
		data.IsLoggedInUser = true
	}
	tplPostUpdate.Execute(w, data)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	postToBeUpdated := r.URL.Query().Get("id")
	body := r.FormValue("body")
	if _, err := db.Db.Query("update Posts set body =$1 where id=$2", body, postToBeUpdated); err != nil {
		w.Write([]byte("<script>alert('Sorry!')</script>"))
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := db.Db.Query("DELETE FROM Posts WHERE id=$1", id)
	if err != nil {
		panic(err.Error())
	}
	log.Println("POST DELETED")
	http.Redirect(w, r, "/", 301)
}
