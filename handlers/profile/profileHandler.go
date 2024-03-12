package profile

import (
	"blog/data"
	"blog/db"
	"blog/handlers/authorization"
	"blog/models"
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func ProfilePage(w http.ResponseWriter, r *http.Request) {
	if len(authorization.Authorized_user) < 1 {
		http.Redirect(w, r, "/login", 302)
	} else {
		tplProfilePage, err := template.ParseFiles("./templates/profile_page.html", "./templates/base.html")
		if err != nil {
			log.Println(err)
		}
		var ProfilePageData data.ProfilePageData
		if len(authorization.Authorized_user) < 1 {
			ProfilePageData.IsLoggedInUser = false
		} else {
			ProfilePageData.IsLoggedInUser = true
		}
		userRows, errQuery := db.Db.Query("SELECT firstname , lastname , email, profile_pic from users where username=$1", authorization.Authorized_user)
		if errQuery != nil {
			log.Println("Query Failed")
		} else {
			defer userRows.Close()
			for userRows.Next() {
				var firstName, lastName, email, profile_pic string
				errQuery := userRows.Scan(&firstName, &lastName, &email, &profile_pic)
				if errQuery != nil {
					log.Println("Sorry")
				}
				if len(firstName) < 1 {
					firstName = " "
				}
				if len(lastName) < 1 {
					lastName = " "
				}
				if len(email) < 1 {
					email = " "
				}
				log.Println(firstName, lastName, email)
				ProfilePageData.FirstName, ProfilePageData.LastName, ProfilePageData.Email, ProfilePageData.ProfilePic = firstName, lastName, email, profile_pic
			}
		}
		log.Println(ProfilePageData.FirstName)
		rows, err := db.Db.Query("select id, body,username from posts where username =$1 order by posts.id desc", authorization.Authorized_user)
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
			post.Id = id
			post.Body = body
			post.Username = username
			post.Editable = true
			post.Deletable = true
			posts = append(posts, post)

		}

		ProfilePageData.Posts = posts
		ProfilePageData.LoggedUser = authorization.Authorized_user
		log.Println(posts)

		//Token generate
		timeNow := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(timeNow, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		ProfilePageData.Token = token
		tplProfilePage.Execute(w, ProfilePageData)
	}

}

func ProfilePageInputHandler(w http.ResponseWriter, r *http.Request) {

	file, handler, err := r.FormFile("uploadfile")
	if handler == nil {
		r.ParseForm()
		firstname := r.FormValue("firstName")
		lastname := r.FormValue("lastName")
		email := r.FormValue("email")
		log.Println(firstname, lastname, email)
		_, errNew := db.Db.Query("update users set firstname=$1, lastname = $2, email=$3 where username =$4", firstname, lastname, email, authorization.Authorized_user)
		if errNew != nil {
			log.Println(errNew)
		} else {
			http.Redirect(w, r, "/profile", 302)
		}
	} else {
		r.ParseMultipartForm(32 << 20)

		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		f, err := os.OpenFile("./images/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		_, errUpdateProfilePic := db.Db.Query("update users set profile_pic =$1 where username =$2", handler.Filename, authorization.Authorized_user)
		if errUpdateProfilePic != nil {
			log.Println("Sorry")
		} else {
			http.Redirect(w, r, "/profile", 302)
		}
	}
}
