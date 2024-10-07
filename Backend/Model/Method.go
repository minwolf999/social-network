package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

/*
This function takes 1 argument:
  - a string who contain a description of the error

The purpose of this function is to Return an error of the application who have make a request to the server.

The function return a string to the user but have no return for the server
*/
func (w *ResponseWriter) Error(err string) {
	time.Sleep(2 * time.Second)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"Error":   http.StatusText(http.StatusUnauthorized),
		"Message": err,
	})
}

// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
//
//  Parse Method for UserData struct
//
// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------

/*
This function takes 1 argument:
  - a map who contain the value of the select and the name of the colum in the db selected

The purpose of this function is to parse the datas into a good structure.

The function return 2 values:
  - an variable of type Auth
  - an error
*/
func (userData *UserData) ParseAuthData() (Auth, error) {
	// We marshal the map to get it in []byte
	serializedData, err := json.Marshal(userData)
	if err != nil {
		return Auth{}, errors.New("internal error: conversion problem")
	}

	// We Unmarshal in the good structure
	var authResult []Auth
	err = json.Unmarshal(serializedData, &authResult)

	return authResult[0], err
}

func (userData *UserData) ParseRegisterData() (Register, error) {
	// We marshal the map to get it in []byte
	serializedData, err := json.Marshal(userData)
	if err != nil {
		return Register{}, errors.New("internal error: conversion problem")
	}

	// We Unmarshal in the good structure
	var registerResult []Register
	err = json.Unmarshal(serializedData, &registerResult)

	return registerResult[0], err
}

func (userData *UserData) ParseCommentsData() ([]Comment, error) {
	// We marshal the map to get it in []byte
	serializedData, err := json.Marshal(userData)
	if err != nil {
		return nil, errors.New("internal error: conversion problem")
	}

	// We Unmarshal in the good structure
	var postResult []Comment
	err = json.Unmarshal(serializedData, &postResult)
	return postResult, err
}

func (userData *UserData) ParseCommentData() (Comment, error) {
	// We marshal the map to get it in []byte
	serializedData, err := json.Marshal(userData)
	if err != nil {
		return Comment{}, errors.New("internal error: conversion problem")
	}

	// We Unmarshal in the good structure
	var postResult Comment
	err = json.Unmarshal(serializedData, &postResult)
	return postResult, err
}

/*
This function takes 1 argument:
  - a array of map who contain the value of the select and the name of the colum in the db selected

The purpose of this function is to parse the datas into a good structure.

The function return 2 values:
  - an variable of type array of Post
  - an error
*/
func (userData *UserData) ParsePostsData() ([]Post, error) {
	var postResult []Post

	for _, v := range *userData {
		var post Post

		// We marshal the map to get it in []byte
		serializedData, err := json.Marshal(v)
		if err != nil {
			return nil, errors.New("internal error: conversion problem")
		}

		// We Unmarshal in the good structure
		if err = json.Unmarshal(serializedData, &post); err != nil {
			return nil, err
		}

		postResult = append(postResult, post)
	}

	return postResult, nil
}

func (userData *UserData) ParsePostData() (Post, error) {
	var post Post

	// We marshal the map to get it in []byte
	serializedData, err := json.Marshal((*userData)[0])
	if err != nil {
		return Post{}, errors.New("internal error: conversion problem")
	}

	// We Unmarshal in the good structure
	if err = json.Unmarshal(serializedData, &post); err != nil {
		return Post{}, err
	}

	return post, nil
}

func (userData *UserData) ParseFollowersData() ([]Follower, error) {
	var res []Follower

	for _, v := range *userData {
		serializedData, err := json.Marshal(v)
		if err != nil {
			return nil, errors.New("internal error: conversion problem")
		}

		var tmp Follower
		if err = json.Unmarshal(serializedData, &tmp); err != nil {
			return nil, err
		}

		res = append(res, tmp)
	}

	return res, nil
}

func (userData *UserData) ParseGroupData() (Group, error) {
	// We marshal the map to get it in []byte
	serializedData, err := json.Marshal(userData)
	if err != nil {
		return Group{}, errors.New("internal error: conversion problem")
	}

	// We Unmarshal in the good structure
	var authResult Group
	err = json.Unmarshal(serializedData, &authResult)

	return Group{}, err
}

// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
//
//	DB Method for Auth struct
//
// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
func (auth *Auth) InsertIntoDb(db *sql.DB) error {
	if auth.Id == "" || auth.Email == "" || auth.Password == "" {
		return errors.New("empty field")
	}

	return InsertIntoDb("Auth", db, auth.Id, auth.Email, auth.Password, auth.ConnectionAttempt)
}

func (auth *Auth) SelectFromDbById(db *sql.DB) error {
	if auth.Id == "" {
		return errors.New("no Id in the struct")
	}

	userData, err := SelectFromDb("Auth", db, map[string]any{"Id": auth.Id})
	if err != nil {
		return err
	}

	*auth, err = userData.ParseAuthData()
	return err
}

func (auth *Auth) SelectFromDbByEmail(db *sql.DB) error {
	if auth.Email == "" {
		return errors.New("no Id in the struct")
	}

	userData, err := SelectFromDb("Auth", db, map[string]any{"Email": auth.Email})
	if err != nil {
		return err
	}

	*auth, err = userData.ParseAuthData()
	return err
}

// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
//
//	DB Method for Register struct
//
// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
func (register *Register) InsertIntoDb(db *sql.DB) error {
	if register.Id == "" || register.Email == "" || register.FirstName == "" || register.LastName == "" || register.BirthDate == "" {
		return errors.New("empty field")
	}

	return InsertIntoDb("UserInfo", db, register.Auth.Id, register.Auth.Email, register.FirstName, register.LastName, register.BirthDate, register.ProfilePicture, register.Username, register.AboutMe)
}

func (register *Register) SelectFromDbById(db *sql.DB) error {
	if register.Id == "" {
		return errors.New("no Id in the struct")
	}

	userData, err := SelectFromDb("UserInfo", db, map[string]any{"Id": register.Id})
	if err != nil {
		return err
	}

	*register, err = userData.ParseRegisterData()
	return err
}

// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
//
//	DB Method for Post struct
//
// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
func (post *Post) InsertIntoDb(db *sql.DB) error {
	if post.Id == "" || post.AuthorId == "" || post.Text == "" || post.CreationDate == "" || (post.Status != "public" && post.Status != "private" && strings.Split(post.Status, " | ")[0] != "almost private") {
		return errors.New("empty field")
	}

	var isGroup = sql.NullString{Valid: false}
	if post.IsGroup != "" {
		isGroup.String = post.IsGroup
		isGroup.Valid = true
	}

	return InsertIntoDb("Post", db, post.Id, post.AuthorId, post.Text, post.Image, post.CreationDate, post.Status, isGroup, 0, 0)
}

func (post *Post) SelectFromDbById(db *sql.DB) error {
	if post.Id == "" {
		return errors.New("no Id in the struct")
	}

	userData, err := SelectFromDb("Post", db, map[string]any{"Id": post.Id})
	if err != nil {
		return err
	}

	// We marshal the map to get it in []byte
	*post, err = userData.ParsePostData()
	return err
}

func (posts *Posts) SelectAllFromDb(db *sql.DB) error {
	userData, err := SelectFromDb("Post", db, map[string]any{})
	if err != nil {
		return err
	}

	// We marshal the map to get it in []byte
	*posts, err = userData.ParsePostsData()
	return err
}

// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
//
//	DB Method for Comment struct
//
// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
func (comment *Comment) InsertIntoDb(db *sql.DB) error {
	if comment.Id == "" || comment.AuthorId == "" || comment.Text == "" || comment.CreationDate == "" || comment.PostId == "" {
		return errors.New("empty field")
	}

	return InsertIntoDb("Comment", db, comment.Id, comment.AuthorId, comment.Text, comment.CreationDate, comment.PostId, 0, 0)
}

func (comment *Comment) SelectFromDbById(db *sql.DB) error {
	if comment.Id == "" {
		return errors.New("no Id in the struct")
	}

	userData, err := SelectFromDb("Comment", db, map[string]any{"Id": comment.Id})
	if err != nil {
		return err
	}

	// We marshal the map to get it in []byte
	*comment, err = userData.ParseCommentData()
	return err
}

func (comments *Comments) SelectAllFromDb(db *sql.DB) error {
	userData, err := SelectFromDb("Post", db, map[string]any{})
	if err != nil {
		return err
	}

	// We marshal the map to get it in []byte
	*comments, err = userData.ParseCommentsData()
	return err
}

// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
//
//	DB Method for Follower struct
//
// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
func (follower *Follower) InsertIntoDb(db *sql.DB) error {
	if follower.Id == "" || follower.UserId == "" || follower.FollowerId == "" {
		return errors.New("empty field")
	}

	return InsertIntoDb("Follower", db, follower.Id, follower.UserId, follower.FollowerId)
}

func (follower *Follower) IsFollowedBy(db *sql.DB) (bool, error) {
	if follower.UserId == "" || follower.FollowerId == "" {
		return false, errors.New("empty id")
	}

	userData, err := SelectFromDb("Follower", db, map[string]any{"UserId": follower.UserId, "FollowedId": follower.FollowerId})
	if err != nil {
		return false, err
	}

	// We marshal the map to get it in []byte
	res, err := userData.ParseFollowersData()
	if err != nil {
		return false, err
	}

	return len(res) == 1, nil
}

func (followers *Followers) SelectFromDbByUserId(db *sql.DB) error {
	if (*followers)[0].UserId == "" {
		return errors.New("no Id in the struct")
	}

	userData, err := SelectFromDb("Follower", db, map[string]any{"UserId": (*followers)[0].UserId})
	if err != nil {
		return err
	}

	*followers, err = userData.ParseFollowersData()
	return err
}

func (followers *Followers) SelectFromDbByFollowerId(db *sql.DB) error {
	if (*followers)[0].FollowerId == "" {
		return errors.New("no Id in the struct")
	}

	userData, err := SelectFromDb("Follower", db, map[string]any{"FollowerId": (*followers)[0].FollowerId})
	if err != nil {
		return err
	}

	// We marshal the map to get it in []byte
	*followers, err = userData.ParseFollowersData()
	return err
}

// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
//
//	DB Method for Group struct
//
// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
func (group *Group) InsertIntoDb(db *sql.DB) error {
	if group.Id == "" || group.LeaderId == "" || group.MemberIds == "" || group.GroupName == "" || group.CreationDate == "" {
		return errors.New("empty field")
	}

	return InsertIntoDb("Groups", db, group.Id, group.LeaderId, group.MemberIds, group.GroupName, group.CreationDate)
}

func (group *Group) SelectFromDbById(db *sql.DB) error {
	if group.Id == "" {
		return errors.New("no Id in the struct")
	}

	userData, err := SelectFromDb("Groups", db, map[string]any{"Id": group.Id})
	if err != nil {
		return err
	}

	*group, err = userData.ParseGroupData()
	return err
}

func (group *Group) SplitMembers() {
	group.SplitMemberIds = strings.Split(group.MemberIds, " | ")
}

func (group *Group) JoinMembers() {
	group.MemberIds = strings.Join(group.SplitMemberIds, " | ")
}
