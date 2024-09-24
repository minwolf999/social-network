package handler

import (
	"database/sql"
	"fmt"
	"log"
	utils "social-network/Utils"

	"golang.org/x/crypto/bcrypt"
)

func ChangeUserName(db *sql.DB, name, uid string) error {
	actualname, err := utils.SelectFromDb("UserInfo", db, map[string]any{"Id": uid})
	if err != nil {
		return err
	}

	userdata, err := utils.ParseRegisterData(actualname[0])
	if err != nil {
		log.Println("Error Parsing Data", err)
		return err
	}

	if userdata.Username == name {
		log.Println("New Username and current Username are the same", err)
		return err
	} else {
		utils.UpdateDb("UserInfos", db, map[string]any{"Username": name}, map[string]any{"Id": uid})
		fmt.Println("Change username Succes")
	}

	return nil
}

func ChangePass(db *sql.DB, newpass, uid string) error {
	actualpass, err := utils.SelectFromDb("UserInfo", db, map[string]any{"Id": uid})
	if err != nil {
		return err
	}

	userdata, err := utils.ParseAuthData(actualpass[0])
	if err != nil {
		log.Println("Error Parsing Data", err)
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userdata.Password), []byte(newpass)); err != nil {
		log.Println("This Password is already used", err)
		return err
	} else {
		utils.UpdateDb("Auth", db, map[string]any{"Password": newpass}, map[string]any{"Id": uid})
		fmt.Println("Change password Succes")
	}
	return nil
}

/*
if err = bcrypt.CompareHashAndPassword([]byte(userdata.Username), []byte(name)); err != nil {
	log.Println("New Username and current Username are the same", err)
	return err
} else { */
