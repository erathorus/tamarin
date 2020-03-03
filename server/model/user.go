package model

type User struct {
	ID             int64  `json:"id"`
	AuthID         string `json:"authId,omitempty"`
	FamilyName     string `json:"familyName"`
	GivenName      string `json:"givenName"`
	ProfilePicture string `json:"profilePicture"`
}
