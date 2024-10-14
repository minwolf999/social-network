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

/*
This function takes 1 argument:
  - a pointer to a UserData object, which contains the data needed to register a user.

The purpose of this function is to parse the registration data into a proper Register structure.

The function returns 2 values:
  - a Register object which contains the parsed registration data
  - an error if something goes wrong during the parsing
*/
func (userData *UserData) ParseRegisterData() (Register, error) {
	// We check if the userData is empty
	if len(*userData) == 0 {
		// Return an error if there is no data
		return Register{}, errors.New("there is no datas")
	}

	// We marshal the userData to convert it to JSON format ([]byte)
	serializedData, err := json.Marshal(userData)
	if err != nil {
		// Return an error if the marshaling fails
		return Register{}, errors.New("internal error: conversion problem")
	}

	// We declare a variable to hold the unmarshaled registration data
	var registerResult []Register

	// We unmarshal the JSON data into the registerResult slice
	err = json.Unmarshal(serializedData, &registerResult)

	// Return the first element of the registerResult slice and any error encountered
	return registerResult[0], err
}

/*
This function takes 1 argument:
  - a pointer to a UserData object, which contains comment data.

The purpose of this function is to parse the comment data into a structured array of Comment objects.

The function returns 2 values:
  - an array of Comment objects
  - an error if something goes wrong during the parsing
*/
func (userData *UserData) ParseCommentsData() ([]Comment, error) {
	// We marshal the userData to convert it to JSON format ([]byte)
	serializedData, err := json.Marshal(userData)
	if err != nil {
		// Return an error if the marshaling fails
		return nil, errors.New("internal error: conversion problem")
	}

	// We declare a variable to hold the unmarshaled comment data
	var postResult []Comment

	// We unmarshal the JSON data into the postResult slice
	err = json.Unmarshal(serializedData, &postResult)

	// Return the result and any error encountered
	return postResult, err
}

/*
This function takes 1 argument:
  - a pointer to a UserData object, which contains comment data.

The purpose of this function is to parse a single comment from the data.

The function returns 2 values:
  - a Comment object, representing the first comment from the parsed data
  - an error if something goes wrong during the parsing
*/
func (userData *UserData) ParseCommentData() (Comment, error) {
	// We call ParseCommentsData to get all comments
	comments, err := userData.ParseCommentsData()

	// We check if there are no comments in the result
	if len(comments) == 0 {
		// Return an error if the data is empty
		return Comment{}, errors.New("there is no data")
	}

	// Return the first comment and any potential error
	return comments[0], err
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

/*
This function takes 1 argument:
  - a pointer to a UserData object, which contains post data.

The purpose of this function is to parse a single post from the data.

The function returns 2 values:
  - a Post object, representing the first post from the parsed data
  - an error if something goes wrong during the parsing
*/
func (userData *UserData) ParsePostData() (Post, error) {
	// We call ParsePostsData to get all posts
	posts, err := userData.ParsePostsData()

	// We check if there are no posts in the result
	if len(posts) == 0 {
		// Return an error if the data is empty
		return Post{}, errors.New("there is no data")
	}

	// Return the first post and any potential error
	return posts[0], err
}

/*
This function takes 1 argument:
  - a pointer to a UserData object, which contains follower data.

The purpose of this function is to parse the follower data into a structured array of Follower objects.

The function returns 2 values:
  - an array of Follower objects
  - an error if something goes wrong during the parsing
*/
func (userData *UserData) ParseFollowersData() ([]Follower, error) {
	// We declare a slice to hold the parsed follower data
	var res []Follower

	// We iterate over each element in the userData
	for _, v := range *userData {
		// We marshal the individual element to convert it to JSON format ([]byte)
		serializedData, err := json.Marshal(v)
		if err != nil {
			// Return an error if the marshaling fails
			return nil, errors.New("internal error: conversion problem")
		}

		// We declare a temporary variable to hold the unmarshaled follower data
		var tmp Follower

		// We unmarshal the JSON data into the tmp variable
		if err = json.Unmarshal(serializedData, &tmp); err != nil {
			// Return the error if the unmarshaling fails
			return nil, err
		}

		// Append the parsed follower data to the result slice
		res = append(res, tmp)
	}

	// Return the result and any error encountered
	return res, nil
}

/*
This function takes 1 argument:
  - a pointer to a UserData object, which contains group data.

The purpose of this function is to parse the group data into a structured Group object.

The function returns 2 values:
  - a Group object, representing the first group from the parsed data
  - an error if something goes wrong during the parsing
*/
func (userData *UserData) ParseGroupData() (Group, error) {
	// We marshal the userData to convert it to JSON format ([]byte)
	serializedData, err := json.Marshal(userData)
	if err != nil {
		// Return an error if the marshaling fails
		return Group{}, errors.New("internal error: conversion problem")
	}

	// We declare a variable to hold the unmarshaled group data
	var groupResult []Group

	// We unmarshal the JSON data into the groupResult slice
	err = json.Unmarshal(serializedData, &groupResult)

	// Return the first group from the groupResult slice and any error encountered
	return groupResult[0], err
}

func (userData *UserData) ParseEventsData() (Events, error) {
	var postResult Events

	for _, v := range *userData {
		var post Event

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

func (userData *UserData) ParseEventData() (Event, error) {
	// We call ParsePostsData to get all posts
	events, err := userData.ParseEventsData()

	// We check if there are no posts in the result
	if len(events) == 0 {
		// Return an error if the data is empty
		return Event{}, errors.New("there is no data")
	}

	// Return the first post and any potential error
	return events[0], err
}

// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
//
//	DB Method for Auth struct
//
// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------

/*
This function takes 2 arguments:
  - a pointer to an Auth object, which contains the authentication data.
  - a pointer to an sql.DB object, representing the database connection.

The purpose of this function is to insert the authentication data into the database.

The function returns 1 value:
  - an error if any field is empty or if the insertion into the database fails
*/
func (auth *Auth) InsertIntoDb(db *sql.DB) error {
	// We check if any of the required fields (Id, Email, Password) are empty
	if auth.Id == "" || auth.Email == "" || auth.Password == "" {
		// Return an error if any field is empty
		return errors.New("empty field")
	}

	// We call InsertIntoDb to insert the authentication data into the "Auth" table in the database
	return InsertIntoDb("Auth", db, auth.Id, auth.Email, auth.Password, auth.ConnectionAttempt)
}

/*
This function takes 2 arguments:
  - a pointer to an Auth object, which will be populated with the authentication data from the database.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any, which contains the conditions (WHERE clause) for selecting the data from the "Auth" table.

The purpose of this function is to retrieve the authentication data from the database based on the given conditions.

The function returns 1 value:
  - an error if the data retrieval or parsing fails
*/
func (auth *Auth) SelectFromDb(db *sql.DB, where map[string]any) error {
	// We call SelectFromDb to retrieve data from the "Auth" table based on the given conditions
	userData, err := SelectFromDb("Auth", db, where)
	if err != nil {
		// Return an error if the data retrieval fails
		return err
	}

	// We parse the retrieved data into the Auth structure and assign it to the auth object
	*auth, err = userData.ParseAuthData()

	// Return any error encountered during parsing
	return err
}

/*
This function takes 3 arguments:
  - a pointer to an Auth object, which represents the authentication data.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any containing the updateData, which holds the values to be updated.
  - a map[string]any containing the where clause, which specifies the conditions for selecting the record(s) to update.

The purpose of this function is to update the authentication data in the "Auth" table based on the provided conditions.

The function returns 1 value:
  - an error if the update operation fails
*/
func (auth *Auth) UpdateDb(db *sql.DB, updateData, where map[string]any) error {
	// We call UpdateDb to update the "Auth" table with the provided data and conditions
	return UpdateDb("Auth", db, updateData, where)
}

/*
This function takes 2 arguments:
  - a pointer to an Auth object, which represents the authentication data to be deleted.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any containing the where clause, which specifies the conditions for selecting the record(s) to delete.

The purpose of this function is to delete authentication data from the "Auth" table based on the provided conditions.

The function returns 1 value:
  - an error if the delete operation fails
*/
func (auth *Auth) DeleteFromDb(db *sql.DB, where map[string]any) error {
	// We call RemoveFromDB to delete the record(s) from the "Auth" table based on the specified conditions
	return RemoveFromDB("Auth", db, where)
}

// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
//
//	DB Method for Register struct
//
// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------

/*
This function takes 1 argument:
  - a pointer to a Register object, which contains the registration data to be inserted into the database.
  - a pointer to an sql.DB object, representing the database connection.

The purpose of this function is to insert the registration data into the "UserInfo" table in the database.

The function returns 1 value:
  - an error if any of the required fields are empty or if the insertion into the database fails
*/
func (register *Register) InsertIntoDb(db *sql.DB) error {
	// We check if any of the required fields (Id, Email, FirstName, LastName, BirthDate) are empty
	if register.Id == "" || register.Email == "" || register.FirstName == "" || register.LastName == "" || register.BirthDate == "" {
		// Return an error if any field is empty
		return errors.New("empty field")
	}

	// We call InsertIntoDb to insert the registration data into the "UserInfo" table in the database
	return InsertIntoDb("UserInfo", db, register.Auth.Id, register.Auth.Email, register.FirstName, register.LastName, register.BirthDate, register.ProfilePicture, register.Username, register.AboutMe, register.GroupsJoined)
}

/*
This function takes 2 arguments:
  - a pointer to a Register object, which will be populated with the registration data retrieved from the database.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any, which contains the conditions (WHERE clause) for selecting the data from the "UserInfo" table.

The purpose of this function is to retrieve the registration data from the database based on the given conditions.

The function returns 1 value:
  - an error if the data retrieval or parsing fails
*/
func (register *Register) SelectFromDb(db *sql.DB, where map[string]any) error {
	// We call SelectFromDb to retrieve data from the "UserInfo" table based on the given conditions
	userData, err := SelectFromDb("UserInfo", db, where)
	if err != nil {
		// Return an error if the data retrieval fails
		return err
	}

	// We parse the retrieved data into the Register structure and assign it to the register object
	*register, err = userData.ParseRegisterData()

	// Return any error encountered during parsing
	return err
}

/*
This function takes 3 arguments:
  - a pointer to a Register object, which represents the registration data to be updated.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any containing the updateData, which holds the values to be updated.
  - a map[string]any containing the where clause, which specifies the conditions for selecting the record(s) to update.

The purpose of this function is to update the registration data in the "UserInfo" table based on the provided conditions.

The function returns 1 value:
  - an error if the update operation fails
*/
func (register *Register) UpdateDb(db *sql.DB, updateData, where map[string]any) error {
	// We call UpdateDb to update the "UserInfo" table with the provided data and conditions
	return UpdateDb("UserInfo", db, updateData, where)
}

/*
This function takes 2 arguments:
  - a pointer to a Register object, which represents the registration data to be deleted.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any containing the where clause, which specifies the conditions for selecting the record(s) to delete.

The purpose of this function is to delete registration data from the "UserInfo" table based on the provided conditions.

The function returns 1 value:
  - an error if the delete operation fails
*/
func (register *Register) DeleteFromDb(db *sql.DB, where map[string]any) error {
	// We call RemoveFromDB to delete the record(s) from the "UserInfo" table based on the specified conditions
	return RemoveFromDB("UserInfo", db, where)
}

/*
This function takes no arguments and is a method of the Register struct.

The purpose of this function is to split the MemberIds string into a slice of individual member IDs, using " | " as the delimiter.

This function updates the SplitGroupsJoined field of the Register struct with the resulting slice.
*/
func (register *Register) SplitGroups() {
	// We split the GroupsJoined string into a slice of strings using " | " as the delimiter
	register.SplitGroupsJoined = strings.Split(register.GroupsJoined, " | ")
}

/*
This function takes no arguments and is a method of the Register struct.

The purpose of this function is to join the SplitGroupsJoined slice into a single string, using " | " as the separator.

This function updates the GroupsJoined field of the Register struct with the resulting string.
*/
func (register *Register) JoinGroups() {
	// We join the SplitMemberIds slice into a single string using " | " as the separator
	register.GroupsJoined = strings.Join(register.SplitGroupsJoined, " | ")
}

// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
//
//	DB Method for Post struct
//
// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------

/*
This function takes 1 argument:
  - a pointer to a Post object, which contains the post data to be inserted into the database.
  - a pointer to an sql.DB object, representing the database connection.

The purpose of this function is to insert the post data into the "Post" table in the database.

The function returns 1 value:
  - an error if any of the required fields are empty or if the insertion into the database fails
*/
func (post *Post) InsertIntoDb(db *sql.DB) error {
	// We check if any of the required fields (Id, AuthorId, Text, CreationDate) are empty
	// We also validate the Status field to ensure it has an acceptable value
	if post.Id == "" || post.AuthorId == "" || post.Text == "" || post.CreationDate == "" ||
		(post.Status != "public" && post.Status != "private" && strings.Split(post.Status, " | ")[0] != "almost private") {
		// Return an error if any field is empty or if Status is invalid
		return errors.New("empty field")
	}

	// We create a sql.NullString to handle the IsGroup field for optional values
	var isGroup = sql.NullString{Valid: false}
	if post.IsGroup != "" {
		isGroup.String = post.IsGroup
		isGroup.Valid = true
	}

	// We call InsertIntoDb to insert the post data into the "Post" table in the database
	return InsertIntoDb("Post", db, post.Id, post.AuthorId, post.Text, post.Image, post.CreationDate, post.Status, isGroup, 0, 0)
}

/*
This function takes 2 arguments:
  - a pointer to a Post object, which will be populated with the post data retrieved from the database.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any, which contains the conditions (WHERE clause) for selecting the data from the "PostDetail" table.

The purpose of this function is to retrieve the post data from the database based on the given conditions.

The function returns 1 value:
  - an error if the data retrieval or parsing fails
*/
func (post *Post) SelectFromDb(db *sql.DB, where map[string]any) error {
	// We call SelectFromDb to retrieve data from the "PostDetail" table based on the given conditions
	userData, err := SelectFromDb("PostDetail", db, where)
	if err != nil {
		// Return an error if the data retrieval fails
		return err
	}

	// We parse the retrieved data into the Post structure and assign it to the post object
	*post, err = userData.ParsePostData()

	// Return any error encountered during parsing
	return err
}

/*
This function takes 3 arguments:
  - a pointer to a Post object, which represents the post data to be updated.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any containing the updateData, which holds the values to be updated.
  - a map[string]any containing the where clause, which specifies the conditions for selecting the record(s) to update.

The purpose of this function is to update the post data in the "Post" table based on the provided conditions.

The function returns 1 value:
  - an error if the update operation fails
*/
func (post *Post) UpdateDb(db *sql.DB, updateData, where map[string]any) error {
	// We call UpdateDb to update the "Post" table with the provided data and conditions
	return UpdateDb("Post", db, updateData, where)
}

/*
This function takes 2 arguments:
  - a pointer to a Post object, which represents the post data to be deleted.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any containing the where clause, which specifies the conditions for selecting the record(s) to delete.

The purpose of this function is to delete post data from the "Post" table based on the provided conditions.

The function returns 1 value:
  - an error if the delete operation fails
*/
func (post *Post) DeleteFromDb(db *sql.DB, where map[string]any) error {
	// We call RemoveFromDB to delete the record(s) from the "Post" table based on the specified conditions
	return RemoveFromDB("Post", db, where)
}

/*
This function takes 2 arguments:
  - a pointer to a Posts object, which will be populated with the post data retrieved from the database.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any, which contains the conditions (WHERE clause) for selecting the data from the "PostDetail" table.

The purpose of this function is to retrieve multiple post data entries from the database based on the given conditions.

The function returns 1 value:
  - an error if the data retrieval or parsing fails
*/
func (post *Posts) SelectFromDb(db *sql.DB, where map[string]any) error {
	// We call SelectFromDb to retrieve data from the "PostDetail" table based on the given conditions
	userData, err := SelectFromDb("PostDetail", db, where)
	if err != nil {
		// Return an error if the data retrieval fails
		return err
	}

	// We parse the retrieved data into the Posts structure and assign it to the post object
	*post, err = userData.ParsePostsData()

	// Return any error encountered during parsing
	return err
}

// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
//
//	DB Method for Comment struct
//
// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------

/*
This function takes 1 argument:
  - a pointer to a Comment object, which contains the comment data to be inserted into the database.
  - a pointer to an sql.DB object, representing the database connection.

The purpose of this function is to insert the comment data into the "Comment" table in the database.

The function returns 1 value:
  - an error if any of the required fields are empty or if the insertion into the database fails
*/
func (comment *Comment) InsertIntoDb(db *sql.DB) error {
	// We check if any of the required fields (Id, AuthorId, Text, CreationDate, PostId) are empty
	if comment.Id == "" || comment.AuthorId == "" || comment.Text == "" || comment.CreationDate == "" || comment.PostId == "" {
		// Return an error if any field is empty
		return errors.New("empty field")
	}

	// We call InsertIntoDb to insert the comment data into the "Comment" table in the database
	return InsertIntoDb("Comment", db, comment.Id, comment.AuthorId, comment.Text, comment.CreationDate, comment.PostId, 0, 0)
}

/*
This function takes 2 arguments:
  - a pointer to a Comment object, which will be populated with the comment data retrieved from the database.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any, which contains the conditions (WHERE clause) for selecting the data from the "CommentDetail" table.

The purpose of this function is to retrieve comment data from the database based on the given conditions.

The function returns 1 value:
  - an error if the Id field is not set, or if the data retrieval or parsing fails
*/
func (comment *Comment) SelectFromDb(db *sql.DB, where map[string]any) error {
	// We check if the Id field in the Comment struct is empty
	if comment.Id == "" {
		// Return an error if no Id is present
		return errors.New("no Id in the struct")
	}

	// We call SelectFromDb to retrieve data from the "CommentDetail" table based on the given conditions
	userData, err := SelectFromDb("CommentDetail", db, where)
	if err != nil {
		// Return an error if the data retrieval fails
		return err
	}

	// We parse the retrieved data into the Comment structure and assign it to the comment object
	*comment, err = userData.ParseCommentData()

	// Return any error encountered during parsing
	return err
}

/*
This function takes 3 arguments:
  - a pointer to a Comment object, which represents the comment data to be updated.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any containing the updateData, which holds the values to be updated.
  - a map[string]any containing the where clause, which specifies the conditions for selecting the record(s) to update.

The purpose of this function is to update the comment data in the "Comment" table based on the provided conditions.

The function returns 1 value:
  - an error if the update operation fails
*/
func (comment *Comment) UpdateDb(db *sql.DB, updateData, where map[string]any) error {
	// We call UpdateDb to update the "Comment" table with the provided data and conditions
	return UpdateDb("Comment", db, updateData, where)
}

/*
This function takes 2 arguments:
  - a pointer to a Comment object, which represents the comment data to be deleted.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any containing the where clause, which specifies the conditions for selecting the record(s) to delete.

The purpose of this function is to delete comment data from the "Comment" table based on the provided conditions.

The function returns 1 value:
  - an error if the delete operation fails
*/
func (comment *Comment) DeleteFromDb(db *sql.DB, where map[string]any) error {
	// We call RemoveFromDB to delete the record(s) from the "Comment" table based on the specified conditions
	return RemoveFromDB("Comment", db, where)
}

/*
This function takes 2 arguments:
  - a pointer to a Comments object, which will be populated with the comment data retrieved from the database.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any, which contains the conditions (WHERE clause) for selecting the data from the "CommentDetail" table.

The purpose of this function is to retrieve multiple comment data entries from the database based on the given conditions.

The function returns 1 value:
  - an error if the data retrieval or parsing fails
*/
func (comments *Comments) SelectFromDb(db *sql.DB, where map[string]any) error {
	// We call SelectFromDb to retrieve data from the "CommentDetail" table based on the given conditions
	userData, err := SelectFromDb("CommentDetail", db, where)
	if err != nil {
		// Return an error if the data retrieval fails
		return err
	}

	// We parse the retrieved data into the Comments structure and assign it to the comments object
	*comments, err = userData.ParseCommentsData()

	// Return any error encountered during parsing
	return err
}

// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
//
//	DB Method for Follower struct
//
// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------

/*
This function takes 1 argument:
  - a pointer to a Follower object, which contains the follower data to be inserted into the database.
  - a pointer to an sql.DB object, representing the database connection.

The purpose of this function is to insert the follower data into the "Follower" table in the database.

The function returns 1 value:
  - an error if any of the required fields are empty or if the insertion into the database fails
*/
func (follower *Follower) InsertIntoDb(db *sql.DB) error {
	// We check if any of the required fields (Id, UserId, FollowerId) are empty
	if follower.Id == "" || follower.UserId == "" || follower.FollowerId == "" {
		// Return an error if any field is empty
		return errors.New("empty field")
	}

	// We call InsertIntoDb to insert the follower data into the "Follower" table in the database
	return InsertIntoDb("Follower", db, follower.Id, follower.UserId, follower.FollowerId)
}

/*
This function takes 1 argument:
  - a pointer to a Follower object, which represents the follower data to check.
  - a pointer to an sql.DB object, representing the database connection.

The purpose of this function is to determine if a specific user is followed by the follower indicated in the Follower object.

The function returns 2 values:
  - a boolean indicating whether the follower is following the specified user.
  - an error if any of the required IDs are empty or if the data retrieval fails.
*/
func (follower *Follower) IsFollowedBy(db *sql.DB) (bool, error) {
	// We check if either the UserId or FollowerId fields are empty
	if follower.UserId == "" || follower.FollowerId == "" {
		// Return false and an error if any field is empty
		return false, errors.New("empty id")
	}

	// We call SelectFromDb to check if there is a record of the follower in the database
	userData, err := SelectFromDb("Follower", db, map[string]any{"UserId": follower.UserId, "FollowedId": follower.FollowerId})
	if err != nil {
		// Return false and the error if the data retrieval fails
		return false, err
	}

	// We parse the retrieved data to check if the follower exists
	res, err := userData.ParseFollowersData()
	if err != nil {
		// Return false and the error if the parsing fails
		return false, err
	}

	// Return true if there is exactly one record of the follower, otherwise return false
	return len(res) == 1, nil
}

/*
This function takes 2 arguments:
  - a pointer to a Followers object, which will be populated with the follower data retrieved from the database.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any containing the conditions (WHERE clause) for selecting the data from the "Follower" table.

The purpose of this function is to retrieve multiple follower data entries from the database based on the given conditions.

The function returns 1 value:
  - an error if the data retrieval or parsing fails
*/
func (followers *Followers) SelectFromDb(db *sql.DB, where map[string]any) error {
	// We call SelectFromDb to retrieve data from the "Follower" table based on the given conditions
	userData, err := SelectFromDb("Follower", db, where)
	if err != nil {
		// Return an error if the data retrieval fails
		return err
	}

	// We parse the retrieved data into the Followers structure and assign it to the followers object
	*followers, err = userData.ParseFollowersData()

	// Return any error encountered during parsing
	return err
}

/*
This function takes 3 arguments:
  - a pointer to a Follower object, which represents the follower data to be updated.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any containing the updateData, which holds the values to be updated.
  - a map[string]any containing the where clause, which specifies the conditions for selecting the record(s) to update.

The purpose of this function is to update the follower data in the "Follower" table based on the provided conditions.

The function returns 1 value:
  - an error if the update operation fails
*/
func (follower *Follower) UpdateDb(db *sql.DB, updateData, where map[string]any) error {
	// We call UpdateDb to update the "Follower" table with the provided data and conditions
	return UpdateDb("Follower", db, updateData, where)
}

/*
This function takes 2 arguments:
  - a pointer to a Follower object, which represents the follower data to be deleted.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any containing the where clause, which specifies the conditions for selecting the record(s) to delete.

The purpose of this function is to delete follower data from the "Follower" table based on the provided conditions.

The function returns 1 value:
  - an error if the delete operation fails
*/
func (follower *Follower) DeleteFromDb(db *sql.DB, where map[string]any) error {
	// We call RemoveFromDB to delete the record(s) from the "Follower" table based on the specified conditions
	return RemoveFromDB("Follower", db, where)
}

// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
//
//	DB Method for Group struct
//
// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------

/*
This function takes 1 argument:
  - a pointer to a Group object, which contains the group data to be inserted into the database.
  - a pointer to an sql.DB object, representing the database connection.

The purpose of this function is to insert the group data into the "Groups" table in the database.

The function returns 1 value:
  - an error if any of the required fields are empty or if the insertion into the database fails
*/
func (group *Group) InsertIntoDb(db *sql.DB) error {
	// We check if any of the required fields (Id, LeaderId, MemberIds, GroupName, CreationDate) are empty
	if group.Id == "" || group.LeaderId == "" || group.MemberIds == "" || group.GroupName == "" || group.CreationDate == "" {
		// Return an error if any field is empty
		return errors.New("empty field")
	}

	// We call InsertIntoDb to insert the group data into the "Groups" table in the database
	return InsertIntoDb("Groups", db, group.Id, group.LeaderId, group.MemberIds, group.GroupName, group.CreationDate)
}

/*
This function takes 2 arguments:
  - a pointer to a Group object, which will be populated with the group data retrieved from the database.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any containing the conditions (WHERE clause) for selecting the data from the "Groups" table.

The purpose of this function is to retrieve group data from the database based on the given conditions.

The function returns 1 value:
  - an error if the data retrieval or parsing fails
*/
func (group *Group) SelectFromDb(db *sql.DB, where map[string]any) error {
	// We call SelectFromDb to retrieve data from the "Groups" table based on the given conditions
	userData, err := SelectFromDb("Groups", db, where)
	if err != nil {
		// Return an error if the data retrieval fails
		return err
	}

	// We parse the retrieved data into the Group structure and assign it to the group object
	*group, err = userData.ParseGroupData()

	// Return any error encountered during parsing
	return err
}

/*
This function takes 3 arguments:
  - a pointer to a Group object, which contains the group data to be updated.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any containing the updateData, which holds the values to be updated.
  - a map[string]any containing the where clause, which specifies the conditions for selecting the record(s) to update.

The purpose of this function is to update the group data in the "Groups" table based on the provided conditions.

The function returns 1 value:
  - an error if the update operation fails
*/
func (group *Group) UpdateDb(db *sql.DB, updateData, where map[string]any) error {
	// We call UpdateDb to update the "Groups" table with the provided data and conditions
	return UpdateDb("Groups", db, updateData, where)
}

/*
This function takes 2 arguments:
  - a pointer to a Group object, which represents the group data to be deleted.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any containing the where clause, which specifies the conditions for selecting the record(s) to delete.

The purpose of this function is to delete group data from the "Groups" table based on the provided conditions.

The function returns 1 value:
  - an error if the delete operation fails
*/
func (group *Group) DeleteFromDb(db *sql.DB, where map[string]any) error {
	// We call RemoveFromDB to delete the record(s) from the "Groups" table based on the specified conditions
	return RemoveFromDB("Groups", db, where)
}

/*
This function takes no arguments and is a method of the Group struct.

The purpose of this function is to split the MemberIds string into a slice of individual member IDs, using " | " as the delimiter.

This function updates the SplitMemberIds field of the Group struct with the resulting slice.
*/
func (group *Group) SplitMembers() {
	// We split the MemberIds string into a slice of strings using " | " as the delimiter
	group.SplitMemberIds = strings.Split(group.MemberIds, " | ")
}

/*
This function takes no arguments and is a method of the Group struct.

The purpose of this function is to join the SplitMemberIds slice into a single string, using " | " as the separator.

This function updates the MemberIds field of the Group struct with the resulting string.
*/
func (group *Group) JoinMembers() {
	// We join the SplitMemberIds slice into a single string using " | " as the separator
	group.MemberIds = strings.Join(group.SplitMemberIds, " | ")
}

// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
//
//	DB Method for Event struct
//
// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------

/*
This function takes 1 argument:
  - a pointer to a Event object, which contains the group data to be inserted into the database.
  - a pointer to an sql.DB object, representing the database connection.

The purpose of this function is to insert the group data into the "Event" table in the database.

The function returns 1 value:
  - an error if any of the required fields are empty or if the insertion into the database fails
*/
func (event *Event) InsertIntoDb(db *sql.DB) error {
	// We check if any of the required fields are empty
	if event.Id == "" || event.OrganisatorId == "" || event.GroupId == "" || event.Title == "" || event.Description == "" || event.DateOfTheEvent == "" {
		// Return an error if any field is empty
		return errors.New("empty field")
	}

	// We call InsertIntoDb to insert the group data into the "Event" table in the database
	return InsertIntoDb("Event", db, event.Id, event.GroupId, event.OrganisatorId, event.Title, event.Description, event.DateOfTheEvent, event.Image)
}

/*
This function takes 2 arguments:
  - a pointer to a Event object, which will be populated with the event data retrieved from the database.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any containing the conditions (WHERE clause) for selecting the data from the "Event" table.

The purpose of this function is to retrieve group data from the database based on the given conditions.

The function returns 1 value:
  - an error if the data retrieval or parsing fails
*/
func (event *Event) SelectFromDb(db *sql.DB, where map[string]any) error {
	// We call SelectFromDb to retrieve data from the "Groups" table based on the given conditions
	userData, err := SelectFromDb("Event", db, where)
	if err != nil {
		// Return an error if the data retrieval fails
		return err
	}

	// We parse the retrieved data into the Group structure and assign it to the group object
	*event, err = userData.ParseEventData()

	// Return any error encountered during parsing
	return err
}

/*
This function takes 3 arguments:
  - a pointer to a Event object, which contains the group data to be updated.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any containing the updateData, which holds the values to be updated.
  - a map[string]any containing the where clause, which specifies the conditions for selecting the record(s) to update.

The purpose of this function is to update the group data in the "Event" table based on the provided conditions.

The function returns 1 value:
  - an error if the update operation fails
*/
func (event *Event) UpdateDb(db *sql.DB, updateData, where map[string]any) error {
	// We call UpdateDb to update the "Groups" table with the provided data and conditions
	return UpdateDb("Event", db, updateData, where)
}

/*
This function takes 2 arguments:
  - a pointer to a Event object, which represents the group data to be deleted.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any containing the where clause, which specifies the conditions for selecting the record(s) to delete.

The purpose of this function is to delete group data from the "Event" table based on the provided conditions.

The function returns 1 value:
  - an error if the delete operation fails
*/
func (event *Event) DeleteFromDb(db *sql.DB, where map[string]any) error {
	// We call RemoveFromDB to delete the record(s) from the "Groups" table based on the specified conditions
	return RemoveFromDB("Event", db, where)
}

// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
//
//	DB Method for JoinEvent struct
//
// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------

/*
This function takes 1 argument:
  - a pointer to a Event object, which contains the group data to be inserted into the database.
  - a pointer to an sql.DB object, representing the database connection.

The purpose of this function is to insert the group data into the "Event" table in the database.

The function returns 1 value:
  - an error if any of the required fields are empty or if the insertion into the database fails
*/
func (joinEvent *JoinEvent) InsertIntoDb(db *sql.DB) error {
	// We check if any of the required fields are empty
	if joinEvent.UserId == "" || joinEvent.EventId == "" {
		// Return an error if any field is empty
		return errors.New("empty field")
	}

	// We call InsertIntoDb to insert the group data into the "Event" table in the database
	return InsertIntoDb("JoinEvent", db, joinEvent.EventId, joinEvent.UserId)
}

func (joinEvent *JoinEvent) DeleteFromDb(db *sql.DB, where map[string]any) error {
	// We call RemoveFromDB to delete the record(s) from the "Groups" table based on the specified conditions
	return RemoveFromDB("Event", db, where)
}

// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
//
//	DB Method for DeclineEvent struct
//
// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------

/*
This function takes 1 argument:
  - a pointer to a Event object, which contains the group data to be inserted into the database.
  - a pointer to an sql.DB object, representing the database connection.

The purpose of this function is to insert the group data into the "Event" table in the database.

The function returns 1 value:
  - an error if any of the required fields are empty or if the insertion into the database fails
*/
func (declineEvent *DeclineEvent) InsertIntoDb(db *sql.DB) error {
	// We check if any of the required fields are empty
	if declineEvent.UserId == "" || declineEvent.EventId == "" {
		// Return an error if any field is empty
		return errors.New("empty field")
	}

	// We call InsertIntoDb to insert the group data into the "Event" table in the database
	return InsertIntoDb("JoinEvent", db, declineEvent.EventId, declineEvent.UserId)
}

func (declineEvent *DeclineEvent) DeleteFromDb(db *sql.DB, where map[string]any) error {
	// We call RemoveFromDB to delete the record(s) from the "Groups" table based on the specified conditions
	return RemoveFromDB("Event", db, where)
}