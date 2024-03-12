package models

type User struct {
	Username                string
	IsFollowedByCurrentUser bool
	IsOwnedThisAccount      bool
}
