package data

import "blog/models"

type IndexData struct {
	Posts          []models.Post
	LoggedUser     string
	IsLoggedInUser bool
}
