package likes

import (
	"blog/db"
	"blog/handlers/authorization"
	"log"
	"net/http"
)

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	log.Println(id)
	_, err := db.Db.Query("INSERT INTO likes(post_id,user_name,is_liked) values($1,$2,'t')", id, authorization.Authorized_user)
	if err != nil {
		log.Println("Inserting Like Failed")
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func UnlikeHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	postLikeByUserIdQuery, err := db.Db.Query("SELECT id from likes where post_id=$1 and user_name =$2", id, authorization.Authorized_user)
	if err != nil {
		log.Println("Failed")
	} else {
		for postLikeByUserIdQuery.Next() {
			var likeId int64
			err = postLikeByUserIdQuery.Scan(&likeId)
			log.Println(likeId)
			_, err = db.Db.Query("update likes set post_id = $1 , user_name = $2,is_liked ='f' where id =$3", id, authorization.Authorized_user, likeId)
			if err != nil {
				log.Println("Post Like Query Failed")
			} else {
				http.Redirect(w, r, "/", 302)
			}
		}
	}

}
