package data

import "blog/models"

type ProfilePageData struct {
	Posts          []models.Post
	FirstName      string
	LastName       string
	Email          string
	ProfilePic     string
	LoggedUser     string
	IsLoggedInUser bool
	Token          string
}
