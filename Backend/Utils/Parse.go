package utils

import (
	"encoding/json"
	"errors"

	model "social-network/Model"
)

/*
This function takes 1 argument:
  - a map who contain the value of the select and the name of the colum in the db selected

The purpose of this function is to parse the datas into a good structure.

The function return 2 values:
  - an variable of type Auth
  - an error
*/
func ParseAuthData(userData map[string]any) (model.Auth, error) {
	// We marshal the map to get it in []byte
	serializedData, err := json.Marshal(userData)
	if err != nil {
		return model.Auth{}, errors.New("internal error: conversion problem")
	}

	// We Unmarshal in the good structure
	var authResult model.Auth
	err = json.Unmarshal(serializedData, &authResult)

	return authResult, err
}

func ParseRegisterData(userData map[string]any) (model.Register, error) {
	// We marshal the map to get it in []byte
	serializedData, err := json.Marshal(userData)
	if err != nil {
		return model.Register{}, errors.New("internal error: conversion problem")
	}

	// We Unmarshal in the good structure
	var registerResult model.Register
	err = json.Unmarshal(serializedData, &registerResult)

	return registerResult, err
}

func ParseCommentData(userData []map[string]any) ([]model.Comment, error) {
	// We marshal the map to get it in []byte
	serializedData, err := json.Marshal(userData)
	if err != nil {
		return nil, errors.New("internal error: conversion problem")
	}

	// We Unmarshal in the good structure
	var postResult []model.Comment
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
func ParsePostData(userData []map[string]any) ([]model.Post, error) {
	// We marshal the map to get it in []byte
	serializedData, err := json.Marshal(userData)
	if err != nil {
		return nil, errors.New("internal error: conversion problem")
	}

	// We Unmarshal in the good structure
	var postResult []model.Post
	err = json.Unmarshal(serializedData, &postResult)
	return postResult, err
}

func ParseFollowerData(follow []map[string]any) (struct{ Follow []model.Follower }, error) {
	var res struct {
		Follow []model.Follower
	}

	for _, v := range follow {
		serializedData, err := json.Marshal(v)
		if err != nil {
			return res, errors.New("internal error: conversion problem")
		}

		var tmp model.Follower
		if err = json.Unmarshal(serializedData, &tmp); err != nil {
			return res, err
		}

		res.Follow = append(res.Follow, tmp)

	}

	return res, nil
}
