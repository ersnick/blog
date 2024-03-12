package authorization

import (
	"blog/db"
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

var Authorized_user string

// Render Login page
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	Authorized_user = ""
	tplLogin, err := template.ParseFiles("./templates/login.html", "./templates/base.html")
	if err != nil {
		log.Println(err)
	}
	tplLogin.Execute(w, nil)
}

// Handling post data from login page
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	result := db.Db.QueryRow("select password from users where username=$1", username)
	var obtainedPassword string
	err := result.Scan(&obtainedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("<script>alert('No user exist!')</script>"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if obtainedPassword != password {
		w.Write([]byte("<script>alert('Login Failed!')</script>"))
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		Authorized_user = username
		http.Redirect(w, r, "/", 302)
	}

}

// Render sign up page
func SignUpPageHandler(w http.ResponseWriter, r *http.Request) {
	tplRegister, err := template.ParseFiles("./templates/signup.html", "./templates/base.html")
	if err != nil {
		log.Println(err)
	}
	tplRegister.Execute(w, nil)
}

// Getting post request and handle them
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password1 := r.FormValue("password1")
	password2 := r.FormValue("password2")
	_username, _password1, _password2 := false, false, false
	_username = !IsEmpty(username)
	_password1 = !IsEmpty(password1)
	_password2 = !IsEmpty(password2)
	if _username && _password1 && _password2 {
		if string(password1) != string(password2) {
			http.Redirect(w, r, "/signup", 302)
		} else {
			if _, err := db.Db.Query("insert into users values ($1, $2)", username, password1); err != nil {
				w.Write([]byte("<script>alert('Error occurred!')</script>"))
			} else {
				w.Write([]byte("<script>alert('Success! Please login')</script>"))
			}
		}
	} else {
		w.Write([]byte("<script>alert('Sorry! Fields can not be empty')</script>"))
	}

}

func Logout(w http.ResponseWriter, r *http.Request) {
	Authorized_user = ""
	http.Redirect(w, r, "/login", 301)

}

// Delete user
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	_, err := db.Db.Query("DELETE FROM users WHERE username=$1", username)
	if err != nil {
		panic(err.Error())
	}
	log.Println("DELETE")
	http.Redirect(w, r, "/", 301)
}

func IsEmpty(data string) bool {
	if len(data) <= 0 {
		return true
	} else {
		return false
	}
}
