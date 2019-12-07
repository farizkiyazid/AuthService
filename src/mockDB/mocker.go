package mockDB

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type User struct {
	ID          string `json:"ID"`
	User       	string `json:"User"`
	Pass 		string `json:"Pass"`
	Token 		string `json:"Token"`
}

type allUsers []User

var Users = allUsers{
	{
		ID:          "1",
		User:       "eldy",
		Pass: 		"boy",
		Token:		"",
	},
}


func CheckOneUser(w http.ResponseWriter, username string, pass string) bool {

	for _, singleUser := range Users {
		if singleUser.User == username {
			if singleUser.Pass == pass {
				// json.NewEncoder(w).Encode(singleUser)
				return true
			}
		}
	}
	return false
}

func CheckUserToken(w http.ResponseWriter, username string) string {

	for _, singleUser := range Users {
		if singleUser.User == username {
			return singleUser.Token
		}
	}
	return ""
}


func SetToken(username string, token string) {
	fmt.Println("Masuk")
	for _, singleUser := range Users {
		if singleUser.User == username {
			singleUser.Token = token
		}
	}
	fmt.Println(Users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the user id and password only in order to log in")
	}
	
	json.Unmarshal(reqBody, &newUser)
	Users = append(Users, newUser)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newUser)
}