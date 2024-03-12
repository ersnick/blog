package data

import "blog/models"

type UpdatePageData struct {
	Post           models.Post
	LoggedUser     string
	IsLoggedInUser bool
}
