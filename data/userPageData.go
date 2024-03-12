package data

import "blog/models"

type UserPageData struct {
	Users          []models.User
	LoggedUser     string
	IsLoggedInUser bool
}
