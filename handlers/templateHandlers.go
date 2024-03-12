package handlers

import (
	"blog/data"
	"blog/db"
	"blog/handlers/authorization"
	"blog/models"
	"html/template"
	"log"
	"net/http"
)

var tplUpdate = template.Must(template.ParseFiles("./templates/update.html"))

// UpdatePage to render user information update page
func UpdatePage(w http.ResponseWriter, r *http.Request) {
	tplUpdate.Execute(w, authorization.Authorized_user)
}

// UpdateHandler for handling submitted update data
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	userToBeUpdated := r.URL.Query().Get("username")
	username := r.FormValue("username")
	password1 := r.FormValue("password1")
	password2 := r.FormValue("password2")
	_username, _password1, _password2 := false, false, false
	_username = !authorization.IsEmpty(username)
	_password1 = !authorization.IsEmpty(password1)
	_password2 = !authorization.IsEmpty(password2)
	if _username && _password1 && _password2 {
		if string(password1) != string(password2) {
			http.Redirect(w, r, "/signup", 302)
		} else {
			if _, err := db.Db.Query("update users set username=$1,password=$2 where username =$3", username, password1, userToBeUpdated); err != nil {
				w.Write([]byte("<script>alert('Error occurred!')</script>"))
			} else {
				authorization.Authorized_user = username
				http.Redirect(w, r, "/", 302)
			}
		}
	} else {
		w.Write([]byte("<script>alert('Sorry! Fields can not be empty')</script>"))
	}

}

// Render and handle Index page
func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Db.Query("select posts.id,posts.body,posts.username from posts")
	var posts []models.Post
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var body string
		var username string
		var post models.Post
		err = rows.Scan(&id, &body, &username)
		if err != nil {
			log.Println(err)
		}
		post.Likes = 0
		likesQuery, err := db.Db.Query("SELECT is_liked from likes where post_id=$1", id)
		if err != nil {
			log.Println("Likes query failed")
		}
		for likesQuery.Next() {
			var isLiked bool
			err = likesQuery.Scan(&isLiked)
			if isLiked {
				post.Likes += 1
			}

		}
		post.IsLikedByCurrentUser = false
		isLikedByCurrentUserQuery, err := db.Db.Query("SELECT is_liked from likes where user_name=$1 and post_id=$2", authorization.Authorized_user, id)
		if err != nil {
			log.Println("Is Liked by current User query failed")
		}
		for isLikedByCurrentUserQuery.Next() {
			var isLiked bool
			isLikedByCurrentUserQuery.Scan(&isLiked)
			if isLiked {
				post.IsLikedByCurrentUser = true
			}
		}
		post.Id = id
		post.Body = body
		post.Username = username
		if authorization.Authorized_user == username {
			post.Editable = true
			post.Deletable = true
		}
		log.Println(post.IsLikedByCurrentUser)
		posts = append(posts, post)

	}
	var IndexPageData data.IndexData
	if len(authorization.Authorized_user) < 1 {
		IndexPageData.IsLoggedInUser = false
	} else {
		IndexPageData.IsLoggedInUser = true
	}
	IndexPageData.Posts = posts
	IndexPageData.LoggedUser = authorization.Authorized_user
	tm := template.Must(template.ParseFiles("./templates/index.html", "./templates/base.html"))
	errortm := tm.Execute(w, IndexPageData)
	log.Println(errortm)
}

// Post request of posted data
func Indexhandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	body := r.FormValue("body")
	log.Println(body)
	user := authorization.Authorized_user
	if _, err := db.Db.Query("insert into Posts(body,username) values ($1, $2)", body, user); err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", 302)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

func UsersPageHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Db.Query("SELECT username FROM users")
	var users []models.User
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		var user models.User
		err = rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		if username == authorization.Authorized_user {
			user.IsOwnedThisAccount = true
		}

		user.Username = username
		users = append(users, user)

	}
	var UserPageData data.UserPageData
	if len(authorization.Authorized_user) < 1 {
		UserPageData.IsLoggedInUser = false
	} else {
		UserPageData.IsLoggedInUser = true
	}
	UserPageData.Users = users
	UserPageData.LoggedUser = authorization.Authorized_user
	log.Println(users)
	tm := template.Must(template.ParseFiles("./templates/users.html", "./templates/base.html"))
	errortm := tm.Execute(w, UserPageData)
	log.Println(errortm)
}
