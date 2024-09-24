package utils

import "testing"

func TestParseAuthData(t *testing.T) {
	testMap := map[string]any{
		"Email":    "unemail@gmail.com",
		"Password": "MonMotDePasse123!",
	}

	userData, err := ParseAuthData(testMap)
	if err != nil {
		t.Errorf("Error during the parse: %v", err)
		return
	}

	if userData.Email != testMap["Email"] {
		t.Errorf("Email before and after the parse are not the same")
		return
	}

	if userData.Password != testMap["Password"] {
		t.Errorf("password before and after the parse are not the same")
		return
	}
}

func TestParseRegisterData(t *testing.T) {
	testMap := map[string]any{
		"FirstName": "Test",
		"LastName":  "Test",
		"Gender":    "non binary",
		"AboutMe":   "I'm a test",
	}

	userData, err := ParseRegisterData(testMap)
	if err != nil {
		t.Errorf("Error during the parse: %v", err)
		return
	}

	if userData.FirstName != testMap["FirstName"] {
		t.Errorf("firstname before and after the parse are not the same")
		return
	}

	if userData.LastName != testMap["LastName"] {
		t.Errorf("lastname before and after the parse are not the same")
		return
	}

	if userData.Gender != testMap["Gender"] {
		t.Errorf("gender before and after the parse are not the same")
		return
	}

	if userData.AboutMe != testMap["AboutMe"] {
		t.Errorf("gender before and after the parse are not the same")
		return
	}
}

func TestParsePostData(t *testing.T) {
	testMap := []map[string]any{
		{
			"AuthorId": "id",
			"Text":     "Hello wold!",
			"IsGroup":  "0",
		},
	}

	userData, err := ParsePostData(testMap)
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

func TestParseCommentData(t *testing.T) {
	testMap := []map[string]any{
		{
			"AuthorId": "id",
			"Text":     "Hello wold!",
		},
	}

	userData, err := ParseCommentData(testMap)
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

func TestParseFollowerData(t *testing.T) {
	testMap := []map[string]any{
		{
			"UserId": "U_Id",
			"FollowerId":     "F_Id",
		},
	}

	userData, err := ParseFollowerData(testMap)
	if err != nil {
		t.Fatalf("Error during the parse: %v", err)
		return
	}

	if userData[0].UserId != testMap[0]["UserId"] {
		t.Fatal("Text before and after the parse are not the same")
		return
	}

	if userData[0].FollowerId != testMap[0]["FollowerId"] {
		t.Fatal("AuthorId before and after the parse are not the same")
		return
	}
}