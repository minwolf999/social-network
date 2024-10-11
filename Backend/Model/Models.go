package model

import (
	"net/http"
)

var (
	SecretKey = "tYrEQins27rw0ehqkKfJE0Ofxyd6r8QISFtpomcIILFUfRacmDuBa3nS9NXTpZfV99E1AEaU"
)

type Auth struct {
	Id                string `json:"Id"`
	Email             string `json:"Email"`
	Password          string `json:"Password"`
	ConfirmPassword   string `json:"ConfirmPassword"`
	ConnectionAttempt int    `json:"ConnectionAttempt"`
}

type Register struct {
	Auth      `json:",inline"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	BirthDate string `json:"BirthDate"`

	// OPTIONNAL
	ProfilePicture string `json:"ProfilePicture"`
	Username       string `json:"Username"`
	AboutMe        string `json:"AboutMe"`
	Gender         string `json:"Gender"`

	GroupsJoined      string `json:"GroupsJoined"`
	SplitGroupsJoined []string
}

type User struct {
	Auth     Auth
	Register Register
}

type Post struct {
	Id           string `json:"Id"`
	AuthorId     string `json:"AuthorId"`
	Register     `json:",inline"`
	Text         string `json:"Text"`
	Image        string `json:"Image"`
	CreationDate string `json:"CreationDate"`
	Status       string `json:"Status"`
	IsGroup      string `json:"IsGroup"`
	LikeCount    int    `json:"LikeCount"`
	DislikeCount int    `json:"DislikeCount"`
}

type Posts []Post

type Comment struct {
	Id           string `json:"Id"`
	AuthorId     string `json:"AuthorId"`
	Text         string `json:"Text"`
	CreationDate string `json:"CreationDate"`
	PostId       string `json:"PostId"`
	LikeCount    int    `json:"LikeCount"`
	DislikeCount int    `json:"DislikeCount"`
	Register     `json:",inline"`
}

type Comments []Comment

type Follower struct {
	Id         string `json:"Id"`
	UserId     string `json:"UserId"`
	FollowerId string `json:"FollowerId"`
}

type Followers []Follower

type Group struct {
	Id             string `json:"Id"`
	LeaderId       string `json:"LeaderId"`
	MemberIds      string `json:"MemberIds"`
	SplitMemberIds []string
	GroupName      string `json:"GroupName"`
	CreationDate   string `json:"CreationDate"`
}

type Event struct {
	Id             string `json:"Id"`
	GroupId        string `json:"GroupId"`
	OrganisatorId  string `json:"OrganisatorId"`
	Title          string `json:"Title"`
	Description    string `json:"Description"`
	DateOfTheEvent string `json:"DateOfTheEvent"`
	Image          string `json:"Image"`
}

type Events []Event

type JoinEvent struct {
	UserId        string `json:"UserId"`
	EventId       string `json:"EventId"`
}

type DeclineEvent struct {
	UserId        string `json:"UserId"`
	EventId       string `json:"EventId"`
}

type ResponseWriter struct {
	http.ResponseWriter
}

type UserData []map[string]any
