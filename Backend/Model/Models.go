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

	ProfilePicture string `json:"ProfilePicture"`
	Username       string `json:"Username"`
	AboutMe        string `json:"AboutMe"`
	Gender         string `json:"Gender"`

	Status string `json:"Status"`

	GroupsJoined      string `json:"GroupsJoined"`
	SplitGroupsJoined []string
}

type Users []Register

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
	Image        string `json:"Image"`
	CreationDate string `json:"CreationDate"`
	PostId       string `json:"PostId"`

	LikeCount    int `json:"LikeCount"`
	DislikeCount int `json:"DislikeCount"`
	Register     `json:",inline"`
}
type Comments []Comment

type Follower struct {
	Id         string `json:"Id"`
	UserId     string `json:"UserId"`
	FollowerId string `json:"FollowerId"`
}
type Followers []Follower

type FollowRequest struct {
	UserId     string `json:"UserId"`
	FollowerId string `json:"FollowerId"`
}
type FollowRequests []FollowRequest

type Group struct {
	Id             string `json:"Id"`
	LeaderId       string `json:"LeaderId"`
	MemberIds      string `json:"MemberIds"`
	SplitMemberIds []string
	GroupName      string `json:"GroupName"`
	CreationDate   string `json:"CreationDate"`
}

type JoinGroupRequest struct {
	UserId string `json:"UserId"`
	GroupId string `json:"GroupId"`
}
type JoinGroupRequests []JoinGroupRequest

type Groups []Group

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

type EventDetail struct {
	Id          string `json:"Id"`
	GroupName   string `json:"GroupName"`
	Organisator string `json:"Organisator"`

	Title          string `json:"Title"`
	Description    string `json:"Description"`
	DateOfTheEvent string `json:"DateOfTheEvent"`

	JoinUsers    string `json:"JoinUsers"`
	JoinUsersTab []string

	DeclineUsers    string `json:"DeclineUsers"`
	DeclineUsersTab []string
}

type JoinEvent struct {
	UserId  string `json:"UserId"`
	EventId string `json:"EventId"`
}
type JoinEvents []JoinEvent

type DeclineEvent JoinEvent
type DeclineEvents []DeclineEvent

type ResponseWriter struct {
	http.ResponseWriter
}

type UserData []map[string]any
