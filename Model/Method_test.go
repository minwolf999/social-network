package model

import (
	"net/http/httptest"
	"testing"
	"time"
)

func TestError(t *testing.T) {
	recorder := httptest.NewRecorder()

	w := &ResponseWriter{recorder}

	start := time.Now()
	w.Error("testing")
	if time.Since(start) < 2*time.Second {
		t.Fatal("The 2 second delai have not been respected")
		return
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("The Content-Type need to be application/json, but is : %v", contentType)
	}
}

// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
//
//	Parse functions
//
// ----------------------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------------------
func TestParseAuthData(t *testing.T) {
	testMap := UserData{
		{
			"Email":    "unemail@gmail.com",
			"Password": "MonMotDePasse123!",
		},
	}

	userData, err := testMap.ParseAuthData()
	if err != nil {
		t.Errorf("Error during the parse: %v", err)
		return
	}

	if userData.Email != testMap[0]["Email"] {
		t.Errorf("Email before and after the parse are not the same")
		return
	}

	if userData.Password != testMap[0]["Password"] {
		t.Errorf("password before and after the parse are not the same")
		return
	}
}

func TestParseRegisterData(t *testing.T) {
	testMap := UserData{
		{
			"FirstName": "Test",
			"LastName":  "Test",
			"Gender":    "non binary",
			"AboutMe":   "I'm a test",
		},
	}

	userData, err := testMap.ParseRegisterData()
	if err != nil {
		t.Errorf("Error during the parse: %v", err)
		return
	}

	if userData.FirstName != testMap[0]["FirstName"] {
		t.Errorf("firstname before and after the parse are not the same")
		return
	}

	if userData.LastName != testMap[0]["LastName"] {
		t.Errorf("lastname before and after the parse are not the same")
		return
	}

	if userData.Gender != testMap[0]["Gender"] {
		t.Errorf("gender before and after the parse are not the same")
		return
	}

	if userData.AboutMe != testMap[0]["AboutMe"] {
		t.Errorf("gender before and after the parse are not the same")
		return
	}
}

func TestParsePostsData(t *testing.T) {
	testMap := UserData{
		{
			"AuthorId": "id",
			"Text":     "Hello wold!",
			"IsGroup":  "0",
		},
	}

	userData, err := testMap.ParsePostsData()
	if err != nil {
		t.Errorf("Error during the parse: %v", err)
		return
	}

	if userData[0].Text != testMap[0]["Text"] {
		t.Errorf("Text before and after the parse are not the same")
		return
	}

	if userData[0].IsGroup != testMap[0]["IsGroup"] {
		t.Errorf("IsGroup before and after the parse are not the same")
		return
	}

	if userData[0].AuthorId != testMap[0]["AuthorId"] {
		t.Errorf("AuthorId before and after the parse are not the same")
		return
	}
}

func TestParsePostData(t *testing.T) {
	testMap := UserData{
		{
			"AuthorId": "id",
			"Text":     "Hello wold!",
			"IsGroup":  "0",
		},
	}

	userData, err := testMap.ParsePostData()
	if err != nil {
		t.Errorf("Error during the parse: %v", err)
		return
	}

	if userData.Text != testMap[0]["Text"] {
		t.Errorf("Text before and after the parse are not the same")
		return
	}

	if userData.IsGroup != testMap[0]["IsGroup"] {
		t.Errorf("IsGroup before and after the parse are not the same")
		return
	}

	if userData.AuthorId != testMap[0]["AuthorId"] {
		t.Errorf("AuthorId before and after the parse are not the same")
		return
	}
}

func TestParseCommentsData(t *testing.T) {
	testMap := UserData{
		{
			"AuthorId": "id",
			"Text":     "Hello wold!",
		},
	}

	userData, err := testMap.ParseCommentsData()
	if err != nil {
		t.Fatalf("Error during the parse: %v", err)
		return
	}

	if userData[0].Text != testMap[0]["Text"] {
		t.Fatal("Text before and after the parse are not the same")
		return
	}

	if userData[0].AuthorId != testMap[0]["AuthorId"] {
		t.Fatal("AuthorId before and after the parse are not the same")
		return
	}
}

func TestParseCommentData(t *testing.T) {
	testMap := UserData{
		{
			"AuthorId": "id",
			"Text":     "Hello wold!",
		},
	}

	userData, err := testMap.ParseCommentData()
	if err != nil {
		t.Fatalf("Error during the parse: %v", err)
		return
	}

	if userData.Text != testMap[0]["Text"] {
		t.Fatal("Text before and after the parse are not the same")
		return
	}

	if userData.AuthorId != testMap[0]["AuthorId"] {
		t.Fatal("AuthorId before and after the parse are not the same")
		return
	}
}

func TestParseFollowerData(t *testing.T) {
	testMap := UserData{
		{
			"FollowerId": "U_Id",
			"FollowedId": "F_Id",
		},
	}

	userData, err := testMap.ParseFollowersData()
	if err != nil {
		t.Fatalf("Error during the parse: %v", err)
		return
	}

	if userData[0].FollowerId != testMap[0]["FollowerId"] {
		t.Fatal("Text before and after the parse are not the same")
		return
	}

	if userData[0].FollowedId != testMap[0]["FollowedId"] {
		t.Fatal("AuthorId before and after the parse are not the same")
		return
	}
}

func TestParseGroupData(t *testing.T) {
	testMap := UserData{
		{
			"Id":        "id",
			"LeaderId":  "Leader",
			"GroupName": "Group",
		},
	}

	userData, err := testMap.ParseGroupData()
	if err != nil {
		t.Fatalf("Error during the parse: %v", err)
		return
	}

	if userData.Id != testMap[0]["Id"] {
		t.Fatal("Text before and after the parse are not the same")
		return
	}

	if userData.LeaderId != testMap[0]["LeaderId"] {
		t.Fatal("Text before and after the parse are not the same")
		return
	}

	if userData.GroupName != testMap[0]["GroupName"] {
		t.Fatal("Text before and after the parse are not the same")
		return
	}
}
