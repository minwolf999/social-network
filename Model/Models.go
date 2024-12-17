package model

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	SecretKey          = "tYrEQins27rw0ehqkKfJE0Ofxyd6r8QISFtpomcIILFUfRacmDuBa3nS9NXTpZfV99E1AEaU"
	ConnectedWebSocket WebSocket
)

type WebSocket struct {
	Conn map[string]*websocket.Conn
	Mu   sync.Mutex
}

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
	Banner         string `json:"Banner"`
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
	Id string `json:"Id"`

	FollowerId        string `json:"FollowerId"`
	Follower_Name     string `json:"Follower_Name"`
	Follower_Username string `json:"Follower_Username"`
	Follower_Picture  string `json:"Follower_Picture"`

	FollowedId        string `json:"FollowedId"`
	Followed_Name     string `json:"Followed_Name"`
	Followed_Username string `json:"Followed_Username"`
	Followed_Picture  string `json:"Followed_Picture"`
}
type Followers []Follower

type FollowRequest struct {
	FollowerId       string `json:"FollowerId"`
	Follower_Name    string `json:"Follower_Name"`
	Follower_Picture string `json:"Follower_Picture"`

	FollowedId       string `json:"FollowedId"`
	Followed_Name    string `json:"Followed_Name"`
	Followed_Picture string `json:"Followed_Picture"`
}
type FollowRequests []FollowRequest

type Group struct {
	Id               string `json:"Id"`
	LeaderId         string `json:"LeaderId"`
	Leader           string `json:"Leader"`
	MemberIds        string `json:"MemberIds"`
	SplitMemberIds   []string
	GroupName        string `json:"GroupName"`
	GroupDescription string `json:"GroupDescription"`
	CreationDate     string `json:"CreationDate"`

	GroupPicture string `json:"GroupPicture"`
	Banner       string `json:"Banner"`

	NotificationQuantity int
}

type JoinGroupRequest struct {
	UserId    string `json:"UserId"`
	User_Name string `json:"User_Name"`

	GroupId   string `json:"GroupId"`
	GroupName string `json:"GroupName"`
}
type JoinGroupRequests []JoinGroupRequest

type InviteGroupRequest struct {
	SenderId    string `json:"SenderId"`
	Sender_Name string `json:"Sender_Name"`

	GroupId    string `json:"GroupId"`
	Group_Name string `json:"GroupName"`

	ReceiverId    string `json:"ReceiverId"`
	Receiver_Name string `json:"Receiver_Name"`
}
type InviteGroupRequests []InviteGroupRequest

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
	GroupId     string `json:"GroupId"`
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

type Notification struct {
	Id          string `json:"Id"`
	UserId      string `json:"UserId"`
	Status      string `json:"Status"`
	Description string `json:"Description"`
	GroupId     string `json:"GroupId"`
	OtherUserId string `json:"OtherUserId"`
}
type Notifications []Notification

type Message struct {
	Id            string `json:"Id"`
	SenderId      string `json:"SenderId"`
	Sender_Name   string `json:"Sender_Name"`
	CreationDate  string `json:"CreationDate"`
	Message       string `json:"Message"`
	Image         string `json:"Image"`
	ReceiverId    string `json:"ReceiverId"`
	Receiver_Name string `json:"Receiver_Name"`
	GroupId       string `json:"GroupId"`
	Group_Name    string `json:"Group_Name"`
}
type Messages []Message

type ResponseWriter struct {
	http.ResponseWriter
}

type UserData []map[string]any
